package auth

import (
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Login(cookieId string, username string, password string) *jerr.JError {
	username = strings.ToLower(username)
	user, err := db.GetUserByUsername(username)
	if err != nil {
		jerror := jerr.Get("username not found", err)
		jerror.SetDisplayMessage("Incorrect username or password")
		return &jerror
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		jerror := jerr.Get("password hash mismatch", err)
		jerror.SetDisplayMessage("Incorrect username or password")
		return &jerror
	}

	session, err := db.GetSession(cookieId)
	if err != nil {
		jerror := jerr.Get("session not found", err)
		jerror.SetDisplayMessage("There was a server side issue and the event has been logged.")
		return &jerror
	}

	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		jerror := jerr.Get("session save failed", err)
		jerror.SetDisplayMessage("There was a server side issue and the event has been logged.")
		return &jerror
	}

	return nil
}
