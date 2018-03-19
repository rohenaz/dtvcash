package db

import (
	"github.com/btcsuite/btcd/txscript"
	"time"
)

type TransactionOut struct {
	Id            uint   `gorm:"primary_key"`
	TransactionId uint
	Value         int64
	PkScript      []byte `gorm:"unique;"`
	LockString    string
	RequiredSigs  uint
	ScriptClass   uint
	Addresses     []*Address
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (t TransactionOut) ValueInBCH() float64 {
	return float64(t.Value) * 1.e-8
}

func (t TransactionOut) GetScriptClass() string {
	return txscript.ScriptClass(t.ScriptClass).String()
}
