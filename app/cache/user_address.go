package cache

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
)

type UserAddress struct {
	PkHash []byte
}

func GetUserAddress(userId uint) (wallet.Address, error) {
	pkHash, err := GetUserPkHash(userId)
	if err != nil {
		return wallet.Address{}, jerr.Get("error getting pk_hash from cache", err)
	}
	return wallet.GetAddressFromPkHash(pkHash), nil
}

func GetUserPkHash(userId uint) ([]byte, error) {
	var userAddress UserAddress
	err := GetItem(getUserAddressName(userId), &userAddress)
	if err == nil {
		return userAddress.PkHash, nil
	}
	if ! IsMissError(err) {
		return nil, jerr.Get("error getting pk hash from cache", err)
	}
	key, err := db.GetKeyForUser(userId)
	if err != nil {
		return nil, jerr.Get("error getting key from db", err)
	}
	err = SetUserAddress(userId, &UserAddress{PkHash: key.PkHash})
	if err != nil {
		return nil, jerr.Get("error saving cache", err)
	}
	return key.PkHash, nil
}

func SetUserAddress(userId uint, userAddress *UserAddress) error {
	err := SetItem(getUserAddressName(userId), userAddress)
	if err != nil {
		return jerr.Get("error setting user address cache", err)
	}
	return nil
}

func getUserAddressName(userId uint) string {
	return fmt.Sprintf("user-address-%d", userId)
}
