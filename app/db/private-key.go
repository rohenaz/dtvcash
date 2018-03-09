package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type PrivateKey struct {
	Id        uint `gorm:"primary_key"`
	Name      string
	UserId    uint
	Value     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PrivateKey) GetPrivateKey() wallet.PrivateKey {
	return wallet.PrivateKey{
		Secret: p.Value,
	}
}

func CreateNewPrivateKey(name string, userId uint) (*PrivateKey, error) {
	privateKey := wallet.GeneratePrivateKey()
	var dbPrivateKey = &PrivateKey{
		Name:   name,
		UserId: userId,
		Value: privateKey.Secret,
	}
	result := save(dbPrivateKey)
	if result.Error != nil {
		return nil, jerr.Get("error saving private key", result.Error)
	}
	return dbPrivateKey, nil
}

func GetPrivateKeysForUser(userId uint) ([]*PrivateKey, error) {
	var privateKeys []*PrivateKey
	err := find(&privateKeys, PrivateKey{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return privateKeys, nil
}
