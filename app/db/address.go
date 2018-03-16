package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Address struct {
	Id            uint `gorm:"primary_key"`
	Address       []byte
	HeightChecked uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a Address) GetAddress() wallet.Address {
	return wallet.Address{

	}
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
	}
	err := find(&address, address)
	if err == nil {
		return &address, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting address", err)
	}
	err = create(&address)
	if err != nil {
		return nil, jerr.Get("error creating address", err)
	}
	return &address, nil
}
