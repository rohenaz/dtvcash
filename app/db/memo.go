package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type MemoTest struct {
	Id        uint   `gorm:"primary_key"`
	TxHash    []byte `gorm:"unique;size:50"`
	Address   []byte
	PkScript  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m MemoTest) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}
