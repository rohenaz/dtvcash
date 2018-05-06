package auth

import (
	"github.com/memocash/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

func IsLoggedIn(cookieId string) bool {
	session, err := db.GetSession(cookieId)
	if err != nil {
		panic(jerr.Get("error getting session from db", err))
		return false
	}
	if session.UserId > 0 && ! session.HasLoggedOut {
		return true
	}
	return false
}

func GetSessionUser(cookieId string) (*db.User, error) {
	session, err := db.GetSession(cookieId)
	if err != nil || session.UserId == 0 || session.HasLoggedOut {
		if err == nil {
			return nil, jerr.New("Unable to get session user")
		}
		return nil, jerr.Get("Error getting session", err)
	}
	user, err := db.GetUserById(session.UserId)
	if err != nil {
		return nil, jerr.Get("Error getting session user", err)
	}
	return user, nil
}
