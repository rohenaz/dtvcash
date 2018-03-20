package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"strconv"
	"time"
)

const (
	BlockTable         = "Block"
	KeyTable           = "Key"
	TxInTable          = "TxIn"
	TxInTxnOutTable    = "TxIn.TxnOut"
	TxInTxnOutTxnTable = "TxIn.TxnOut.Transaction"
	TxOutTable         = "TxOut"
	TxOutTxnInTable    = "TxOut.TxnIn"
	TxOutTxnInTxnTable = "TxOut.TxnIn.Transaction"
	TxOutAddressTable  = "TxOut.Addresses"
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

func (t *Transaction) GetBlockHeight() string {
	if t.Block == nil {
		return "Unknown"
	}
	return strconv.Itoa(int(t.Block.Height))
}

func (t *Transaction) GetBlockTime() string {
	if t.Block == nil {
		return "-"
	}
	return t.Block.Timestamp.Format("2006-01-02 15:04")
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
	return getTransaction(Transaction{
		Id: transactionId,
	})
}

func getTransaction(txn Transaction) (*Transaction, error) {
	var transaction Transaction
	err := findPreloadColumns([]string{
		BlockTable,
		KeyTable,
		TxInTable,
		TxInTxnOutTable,
		TxInTxnOutTxnTable,
		TxOutTable,
		TxOutTxnInTable,
		TxOutTxnInTxnTable,
		TxOutAddressTable,
	}, &transaction, txn)
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
		TxOutTable,
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
	for index, in := range msg.TxIn {
		var witnesses []*Witness
		for _, witness := range in.Witness {
			witnesses = append(witnesses, &Witness{
				Data: witness,
			})
		}
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if err != nil {
			jerr.Get("error disassembling unlockScript: %s\n", err).Print()
			return nil
		}
		var transactionIn = TransactionIn{
			Index:                 uint(index),
			PreviousOutPointHash:  in.PreviousOutPoint.Hash.CloneBytes(),
			PreviousOutPointIndex: in.PreviousOutPoint.Index,
			SignatureScript:       in.SignatureScript,
			Witnesses:             witnesses,
			Sequence:              in.Sequence,
			UnlockString:          unlockScript,
		}
		transaction.TxIn = append(transaction.TxIn, &transactionIn)
	}
	for index, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			jerr.Get("rror disassembling lockScript: %s\n", err).Print()
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
			Index:        uint(index),
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
