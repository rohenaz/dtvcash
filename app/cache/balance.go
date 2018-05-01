package cache

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

func GetBalance(pkHash []byte) (int64, error) {
	var bal int64
	err := GetItem(getBalanceName(pkHash), &bal)
	if err != nil {
		return 0, jerr.Get("error getting balance", err)
	}
	return bal, nil
}

func SetBalance(pkHash []byte, balance int64) error {
	err := SetItem(getBalanceName(pkHash), balance)
	if err != nil {
		return jerr.Get("error setting balance", err)
	}
	return nil
}

func ClearBalance(pkHash []byte) error {
	err := DeleteItem(getBalanceName(pkHash))
	if err != nil {
		return jerr.Get("error clearing balance", err)
	}
	return nil
}

func getBalanceName(pkHash []byte) string {
	return fmt.Sprintf("balance-%x", pkHash)
}
