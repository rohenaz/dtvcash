package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Address struct {
	Id            uint `gorm:"primary_key"`
	Address       []byte
	KeyId         uint
	Key           *Key
	HeightChecked uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a Address) Save() error {
	result := save(a)
	if result.Error != nil {
		return jerr.Get("error saving address", result.Error)
	}
	return nil
}

func GetAddress(key *Key) (*Address, error) {
	var address = Address{
		Address: key.GetPublicKey().GetAddress().GetScriptAddress(),
		KeyId:   key.Id,
	}
	err := find(&address, address)
	if err == nil {
		return &address, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting address", err)
	}
	address.Key = key
	err = create(&address)
	if err != nil {
		return nil, jerr.Get("error creating address", err)
	}
	return &address, nil
}
