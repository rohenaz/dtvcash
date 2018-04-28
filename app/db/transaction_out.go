package db

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"strings"
	"time"
)

var transactionOutColumns = []string{
	KeyTable,
	TransactionTable,
	TransactionBlockTbl,
}

type TransactionOut struct {
	Id              uint           `gorm:"primary_key"`
	Index           uint32         `gorm:"unique_index:transaction_out_index;"`
	HashString      string
	TransactionHash []byte         `gorm:"unique_index:transaction_out_index;"`
	Transaction     *Transaction   `gorm:"foreignkey:TransactionHash"`
	KeyPkHash       []byte         `gorm:"index:pk_hash"`
	Key             *Key           `gorm:"foreignkey:KeyPkHash"`
	Value           int64
	PkScript        []byte
	LockString      string
	RequiredSigs    uint
	ScriptClass     uint
	TxnInHashString string
	TxnIn           *TransactionIn `gorm:"foreignkey:TxnInHashString"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func getHashString(txHash []byte, index uint32) string {
	hash, err := chainhash.NewHash(txHash)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("out:%s:%d", hash.String(), index)
}

func (t TransactionOut) GetHashString() string {
	return getHashString(t.TransactionHash, t.Index)
}

func (t TransactionOut) IsMemo() bool {
	return strings.HasPrefix(t.LockString, "OP_RETURN 6d")
}

func (t TransactionOut) GetPkHash() []byte {
	split := strings.Split(t.LockString, " ")
	if len(split) != 5 {
		return []byte{}
	}
	pubKey, err := hex.DecodeString(split[2])
	if err != nil {
		return []byte{}
	}
	return pubKey
}

func (t TransactionOut) GetAddressString() string {
	addressPkHash, err := btcutil.NewAddressPubKeyHash(t.KeyPkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error parsing address", err).Print()
		return ""
	}
	return addressPkHash.String()
}

func (t TransactionOut) Save() error {
	result := save(&t)
	if result.Error != nil {
		return jerr.Get("error saving transaction output", result.Error)
	}
	return nil
}

func (t TransactionOut) ValueInBCH() float64 {
	return float64(t.Value) * 1.e-8
}

func (t TransactionOut) HasIn() bool {
	return len(t.TxnInHashString) > 0
}

func (t TransactionOut) IsSpendable() bool {
	if len(t.TxnInHashString) > 0 {
		txIn, _ := GetTransactionInputByHashString(t.TxnInHashString)
		if txIn.Transaction.BlockId > 0 {
			return false
		}
	}
	return true
}

func (t TransactionOut) GetScriptClass() string {
	return txscript.ScriptClass(t.ScriptClass).String()
}

func (t TransactionOut) GetMessage() string {
	if txscript.ScriptClass(t.ScriptClass) == txscript.NullDataTy {
		data, err := txscript.PushedData(t.PkScript)
		if err != nil || len(data) == 0 {
			return ""
		}
		return string(data[0])
	}
	return html.EscapeString(script.GetScriptString(t.PkScript))
}

func GetTransactionOutputById(id uint) (*TransactionOut, error) {
	var txOut TransactionOut
	err := findPreloadColumns(transactionOutColumns, &txOut, TransactionOut{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error finding transaction out", err)
	}
	return &txOut, nil
}

func GetTransactionOutputsForPkHash(pkHash []byte) ([]*TransactionOut, error) {
	var transactionOuts []*TransactionOut
	err := findPreloadColumns([]string{TransactionTable}, &transactionOuts, TransactionOut{
		KeyPkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error finding transaction outputs", err)
	}
	return transactionOuts, nil
}

func GetSpendableTxOut(pkHash []byte, fee int64) (*TransactionOut, error) {
	transactions, err := GetTransactionsForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting transactions", err)
	}
	var txOut *TransactionOut
	for _, txn := range transactions {
		for _, out := range txn.TxOut {
			if out.TxnInHashString == "" && out.Value > fee && bytes.Equal(out.KeyPkHash, pkHash) {
				txOut = out
			}
		}
	}
	if txOut == nil {
		return nil, jerr.New("unable to find an output to spend")
	}
	return txOut, nil
}

func HasSpendable(pkHash []byte) (bool, error) {
	transactions, err := GetTransactionsForPkHash(pkHash)
	if err != nil {
		return false, jerr.Get("error getting transactions", err)
	}
	var txOut *TransactionOut
	for _, txn := range transactions {
		for _, out := range txn.TxOut {
			if out.TxnInHashString == "" && out.Value > 1000 && bytes.Equal(out.KeyPkHash, pkHash) {
				txOut = out
			}
		}
	}
	if txOut == nil {
		return false, nil
	}
	return true, nil
}
