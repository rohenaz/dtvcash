package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/cryptography"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type PrivateKey struct {
	Id        uint `gorm:"primary_key"`
	Name      string
	UserId    uint
	Value     []byte
	PublicKey []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PrivateKey) GetPrivateKey(password string) (*wallet.PrivateKey, error) {
	key, err := cryptography.GenerateKeyFromPassword(password)
	if err != nil {
		return nil, jerr.Get("error generating key from password", err)
	}
	decrypted, err := cryptography.Decrypt(p.Value, key)
	if err != nil {
		return nil, jerr.Get("failed to decrypt", err)
	}
	return &wallet.PrivateKey{
		Secret: decrypted,
	}, nil
}

func (p PrivateKey) GetPublicKey() wallet.PublicKey {
	return wallet.GetPublicKey(p.PublicKey)
}

func CreateNewPrivateKey(name string, password string, userId uint) (*PrivateKey, error) {
	key, err := cryptography.GenerateKeyFromPassword(password)
	if err != nil {
		return nil, jerr.Get("error generating key from password", err)
	}
	privateKey := wallet.GeneratePrivateKey()
	encryptedSecret, err := cryptography.Encrypt(privateKey.Secret, key)
	if err != nil {
		return nil, jerr.Get("failed to encrypt", err)
	}
	var dbPrivateKey = &PrivateKey{
		Name:      name,
		UserId:    userId,
		Value:     encryptedSecret,
		PublicKey: privateKey.GetPublicKey().GetSerialized(),
	}
	result := save(dbPrivateKey)
	if result.Error != nil {
		return nil, jerr.Get("error saving private key", result.Error)
	}
	return dbPrivateKey, nil
}

func GetPublicKeysForUser(userId uint) ([]*PrivateKey, error) {
	var privateKeys []*PrivateKey
	err := find(&privateKeys, PrivateKey{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return privateKeys, nil
}
