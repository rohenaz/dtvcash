package db

import (
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"strconv"
	"time"
)

var transactionColumns = []string{
	BlockTable,
	TxInTable,
	TxInTxnOutTable,
	TxInKeyTable,
	TxInTxnOutTxnTable,
	TxOutTable,
	TxOutKeyTable,
	TxOutTxnInTable,
	TxOutTxnInTxnTable,
}

type Transaction struct {
	Id        uint              `gorm:"primary_key"`
	BlockId   uint
	Block     *Block
	Hash      []byte            `gorm:"unique;"`
	Version   int32
	TxIn      []*TransactionIn  `gorm:"foreignkey:TransactionHash"`
	TxOut     []*TransactionOut `gorm:"foreignkey:TransactionHash"`
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

func (t *Transaction) GetKeyId() uint {
	if len(t.TxIn) != 1 || t.TxIn[0].Key == nil {
		return 0
	}
	return t.TxIn[0].Key.Id
}

type Value struct {
	amount int64
}

func (v Value) GetValue() int64 {
	return v.amount
}

func (v Value) GetValueBCH() float64 {
	return float64(v.GetValue()) * 1.e-8
}

func (t *Transaction) GetValues() map[string]*Value {
	var values = make(map[string]*Value)
	for _, in := range t.TxIn {
		if in.TxnOut == nil || in.Key == nil {
			continue
		}
		mapKey := in.Key.GetAddress().GetEncoded()
		_, ok := values[mapKey]
		if !ok {
			values[mapKey] = &Value{}
		}
		values[mapKey].amount -= in.TxnOut.Value
	}
	for _, out := range t.TxOut {
		if out.Key == nil {
			continue
		}
		mapKey := out.Key.GetAddress().GetEncoded()
		_, ok := values[mapKey]
		if !ok {
			values[mapKey] = &Value{}
		}
		values[mapKey].amount += out.Value
	}
	return values
}

func (t *Transaction) HasFee() bool {
	return t.GetFee() > 0
}

func (t *Transaction) GetFeeBCH() float64 {
	return float64(t.GetFee()) * 1.e-8
}

func (t *Transaction) GetFee() int64 {
	var inputTotal int64
	var outputTotal int64
	for _, in := range t.TxIn {
		if in.TxnOut == nil {
			// Unknown input, unable to calculate fee
			return 0
		}
		inputTotal += in.TxnOut.Value
	}
	for _, out := range t.TxOut {
		outputTotal += out.Value
	}
	return inputTotal - outputTotal
}

func (t *Transaction) Save() error {
	if t.Id == 0 {
		txn, err := GetTransactionByHash(t.Hash)
		if err != nil && ! IsRecordNotFoundError(err) {
			return jerr.Get("error getting transaction by hash", err)
		}
		if txn != nil {
			return jerr.Get("transaction already exists", alreadyExistsError)
		}
	}
	//fmt.Printf("Saving transaction: %#v\n", t)
	result := save(t)
	if result.Error != nil {
		return jerr.Get("error saving transaction", result.Error)
	}
	return nil
}

func (t *Transaction) Delete() error {
	for _, in := range t.TxIn {
		err := in.Delete()
		if err != nil {
			return jerr.Get("error removing transaction input", err)
		}
	}
	for _, out := range t.TxOut {
		result := remove(out)
		if result.Error != nil {
			return jerr.Get("error removing transaction output", result.Error)
		}
	}
	result := remove(t)
	if result.Error != nil {
		return jerr.Get("error removing transaction", result.Error)
	}
	return nil
}

func (t *Transaction) GetChainHash() *chainhash.Hash {
	hash, _ := chainhash.NewHash(t.Hash)
	return hash
}

func (t *Transaction) HasUserKey(userId uint) bool {
	for _, in := range t.TxIn {
		if in.Key != nil && in.Key.UserId == userId {
			return true
		}
	}
	for _, out := range t.TxOut {
		if out.Key != nil && out.Key.UserId == userId {
			return true
		}
	}
	return false
}

func GetTransactionById(transactionId uint) (*Transaction, error) {
	return getTransaction(Transaction{
		Id: transactionId,
	})
}

func getTransaction(whereTxn Transaction) (*Transaction, error) {
	var txn Transaction
	err := findPreloadColumns(transactionColumns, &txn, whereTxn)
	if err != nil {
		return nil, jerr.Get("error finding transaction", err)
	}
	return &txn, nil
}

func GetTransactions() ([]*Transaction, error) {
	var transactions []*Transaction
	err := findPreloadColumns(transactionColumns, &transactions)
	if err != nil {
		return nil, jerr.Get("error finding transactions", err)
	}
	return transactions, nil
}

func GetTransactionsForPkHash(pkHash []byte) ([]*Transaction, error) {
	ins, err := GetTransactionInputsForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting ins", err)
	}
	outs, err := GetTransactionOutputsForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting outs", err)
	}
	var hashes [][]byte
	for _, in := range ins {
		hashes = append(hashes, in.TransactionHash)
	}
	for _, out := range outs {
		hashes = append(hashes, out.TransactionHash)
	}
	query, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query = query.Preload(TxOutTable)
	var transactions []*Transaction
	result := query.Where("hash in (?)", hashes).Find(&transactions)

	if result.Error != nil {
		return nil, jerr.Get("error finding transactions", result.Error)
	}
	return transactions, nil
}

func GetTransactionByHash(hash []byte) (*Transaction, error) {
	var txn = Transaction{
		Hash: hash,
	}
	err := find(&txn, txn)
	if err != nil {
		return nil, jerr.Get("error finding transaction", err)
	}
	return &txn, nil
}

func GetTransactionByHashWithOutputs(hash []byte) (*Transaction, error) {
	var txn = Transaction{
		Hash: hash,
	}
	err := findPreloadColumns([]string{TxOutTable}, &txn, txn)
	if err != nil {
		return nil, jerr.Get("error finding transaction", err)
	}
	return &txn, nil
}

func ConvertMsgToTransaction(msg *wire.MsgTx) (*Transaction, error) {
	txHash := msg.TxHash()
	var txn = Transaction{
		Hash:     txHash.CloneBytes(),
		Version:  msg.Version,
		LockTime: msg.LockTime,
	}
	for index, in := range msg.TxIn {
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if err != nil {
			return nil, jerr.Getf(err, "error disassembling unlockScript: %s", unlockScript)
		}
		var transactionIn = TransactionIn{
			Index:                 uint(index),
			PreviousOutPointHash:  in.PreviousOutPoint.Hash.CloneBytes(),
			PreviousOutPointIndex: in.PreviousOutPoint.Index,
			SignatureScript:       in.SignatureScript,
			Sequence:              in.Sequence,
			UnlockString:          unlockScript,
			TransactionHash:       txHash.CloneBytes(),
		}
		transactionIn.KeyPkHash = transactionIn.GetPkHash()
		transactionIn.HashString = transactionIn.GetHashString()
		transactionIn.TxnOutHashString = getHashString(transactionIn.PreviousOutPointHash, transactionIn.PreviousOutPointIndex)
		txn.TxIn = append(txn.TxIn, &transactionIn)
	}
	for index, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			return nil, jerr.Get("error disassembling lockScript: %s\n", err)
		}
		scriptClass, _, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParamsOld)
		var transactionOut = TransactionOut{
			Index:           uint32(index),
			Value:           out.Value,
			PkScript:        out.PkScript,
			LockString:      lockScript,
			RequiredSigs:    uint(sigCount),
			ScriptClass:     uint(scriptClass),
			TransactionHash: txHash.CloneBytes(),
		}
		if transactionOut.IsMemo() && len(txn.TxIn) == 1 {
			transactionOut.KeyPkHash = txn.TxIn[0].KeyPkHash
		} else {
			transactionOut.KeyPkHash = transactionOut.GetPkHash()
		}
		transactionOut.HashString = transactionOut.GetHashString()
		txn.TxOut = append(txn.TxOut, &transactionOut)
	}
	return &txn, nil
}
