package db

import "time"

type Witness struct {
	Id              uint   `gorm:"primary_key"`
	TransactionInId uint
	Data            []byte `gorm:"unique;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
