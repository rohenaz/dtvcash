package db

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Transaction struct {
	Id        uint   `gorm:"primary_key"`
	KeyId     uint
	Key       *Key
	BlockId   uint
	Block     *Block
	Hash      []byte `gorm:"unique;"`
	Version   int32
	TxIn      []*TransactionIn
	TxOut     []*TransactionOut
	LockTime  uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Transaction) Save() error {
	if t.Id == 0 {
		transaction, err := GetTransactionByHash(t.Hash)
		if err != nil && ! IsRecordNotFoundError(err) {
			return jerr.Get("error getting transaction by hash", err)
		}
		if transaction != nil {
			return jerr.Get("transaction already exists", alreadyExistsError)
		}
	}
	result := save(t)
	if result.Error != nil {
		return jerr.Get("error saving transaction", result.Error)
	}
	return nil
}

func (t *Transaction) GetChainHash() *chainhash.Hash {
	hash, _ := chainhash.NewHash(t.Hash)
	return hash
}

func GetTransactionById(transactionId uint) (*Transaction, error) {
	var transaction Transaction
	err := findPreloadColumns([]string{
		BlockTable,
		KeyTable,
	}, &transaction, Transaction{
		Id: transactionId,
	})
	if err != nil {
		return nil, jerr.Get("error finding transaction", err)
	}
	return &transaction, nil
}

func GetTransactionsForKey(keyId uint) ([]*Transaction, error) {
	var transactions []*Transaction
	err := findPreloadColumns([]string{
		BlockTable,
		KeyTable,
	}, &transactions, Transaction{
		KeyId: keyId,
	})
	if err != nil {
		return nil, jerr.Get("error finding transactions", err)
	}
	return transactions, nil
}

func GetTransactionByHash(hash []byte) (*Transaction, error) {
	var transaction = Transaction{
		Hash: hash,
	}
	err := find(&transaction, transaction)
	if err != nil {
		return nil, jerr.Get("error finding transaction", err)
	}
	return &transaction, nil
}

func ConvertMsgToTransaction(msg *wire.MsgTx) *Transaction {
	txHash := msg.TxHash()
	var transaction = Transaction{
		Hash:     txHash.CloneBytes(),
		Version:  msg.Version,
		LockTime: msg.LockTime,
	}
	for _, in := range msg.TxIn {
		var witnesses []*Witness
		for _, witness := range in.Witness {
			witnesses = append(witnesses, &Witness{
				Data: witness,
			})
		}
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if err != nil {
			fmt.Printf("Error disassembling unlockScript: %s\n", err.Error())
			return nil
		}
		var transactionIn = TransactionIn{
			PreviousOutPointHash:  in.PreviousOutPoint.Hash.CloneBytes(),
			PreviousOutPointIndex: in.PreviousOutPoint.Index,
			SignatureScript:       in.SignatureScript,
			Witnesses:             witnesses,
			Sequence:              in.Sequence,
			UnlockString:          unlockScript,
		}
		transaction.TxIn = append(transaction.TxIn, &transactionIn)
	}
	for _, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			fmt.Printf("Error disassembling lockScript: %s\n", err.Error())
			return nil
		}
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParams)
		var dbAddresses []*Address
		for _, address := range addresses {
			dbAddresses = append(dbAddresses, &Address{
				Data:   address.ScriptAddress(),
				String: address.String(),
			})
		}
		var transactionOut = TransactionOut{
			Value:        out.Value,
			PkScript:     out.PkScript,
			LockString:   lockScript,
			RequiredSigs: uint(sigCount),
			Addresses:    dbAddresses,
			ScriptClass:  uint(scriptClass),
		}
		transaction.TxOut = append(transaction.TxOut, &transactionOut)
	}
	return &transaction
}
