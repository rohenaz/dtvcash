package auth

import (
	"github.com/memocash/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/bcrypt"
)

func UpdatePassword(userId uint, oldPassword string, newPassword string) error {
	user, err := db.GetUserById(userId)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	err = user.Save()
	if err != nil {
		return jerr.Get("error saving password", err)
	}

	return nil
}
