package db

import "time"

type Address struct {
	Id               uint   `gorm:"primary_key"`
	TransactionOutId uint
	Data             []byte `gorm:"unique;"`
	String           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
