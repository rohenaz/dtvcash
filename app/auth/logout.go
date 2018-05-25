package auth

import (
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/jchavannes/jgo/jerr"
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
