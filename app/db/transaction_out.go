package db

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type TransactionOut struct {
	Id            uint   `gorm:"primary_key"`
	Index         uint
	TransactionId uint   `gorm:"unique_index:transaction_out_script;"`
	Transaction   *Transaction
	Value         int64
	PkScript      []byte `gorm:"unique_index:transaction_out_script;"`
	LockString    string
	RequiredSigs  uint
	ScriptClass   uint
	Addresses     []*Address
	TxnInId       uint
	TxnIn         *TransactionIn
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	return t.TxnInId > 0
}

func (t TransactionOut) GetScriptClass() string {
	return txscript.ScriptClass(t.ScriptClass).String()
}
