package auth

import (
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Signup(cookieId string, username string, password string) error {
	username = strings.ToLower(username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := db.CreateUser(username, string(hashedPassword))
	if err != nil {
		return jerr.Get("error signing up", err)
	}
	session, err := db.GetSession(cookieId)
	if err != nil {
		return jerr.Get("error getting session", err)
	}
	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return jerr.Get("error saving session", err)
	}
	return nil
}
