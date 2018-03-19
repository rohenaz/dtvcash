package db

import "time"

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
