package db

import "time"

type Transaction struct {
	Id            uint `gorm:"primary_key"`
	Address       []byte
	KeyId         uint
	Key           *Key
	HeightChecked uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
