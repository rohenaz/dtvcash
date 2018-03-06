package auth

import (
	"github.com/jchavannes/jgo/jerr"
	"git.jasonc.me/main/memo/app/db"
)

func Logout(cookieId string) error {
	session, err := db.GetSession(cookieId)
	if err != nil {
		return jerr.Get("Error getting session", err)
	}

	session.HasLoggedOut = true
	err = session.Save()
	if err != nil {
		return jerr.Get("Error saving session", err)
	}

	return nil
}
