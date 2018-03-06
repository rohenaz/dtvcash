package auth

import (
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/bcrypt"
)

func Login(cookieId string, username string, password string) error {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	session, err := db.GetSession(cookieId)
	if err != nil {
		return jerr.Get("Error getting session", err)
	}

	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return jerr.Get("Error saving session", err)
	}

	return nil
}
