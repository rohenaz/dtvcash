package db

import "time"

type Address struct {
	Id               uint   `gorm:"primary_key"`
	TransactionOutId uint   `gorm:"unique_index:transaction_out_address;"`
	Data             []byte `gorm:"unique_index:transaction_out_address;"`
	String           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
