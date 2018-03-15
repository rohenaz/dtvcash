package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Transaction struct {
	Id            uint `gorm:"primary_key"`
	Address       []byte
	KeyId         uint
	Key           *Key
	HeightChecked uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (t *Transaction) Save() error {
	result := save(t)
	if result.Error != nil {
		return jerr.Get("error saving transaction", result.Error)
	}
	return nil
}
