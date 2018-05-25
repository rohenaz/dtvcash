package auth

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/rohenaz/dtvcash/app/db"
	"golang.org/x/crypto/bcrypt"
)

func DeleteAccount(userId uint, password string) error {
	user, err := db.GetUserById(userId)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	err = user.Delete()
	if err != nil {
		return jerr.Get("error deleting user", err)
	}

	return nil
}
