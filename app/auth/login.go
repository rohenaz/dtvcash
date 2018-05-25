package auth

import (
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	MsgUsernameNotFound = "username not found"
	MsgPasswordMismatch = "password hash mismatch"
)

func IsBadUsernamePasswordError(err error) bool {
	return jerr.HasError(err, MsgUsernameNotFound) || jerr.HasError(err, MsgPasswordMismatch)
}

// Reasonable assumption here but error creating user might not mean already exists in edge cases.
func UserAlreadyExists(err error) bool {
	return jerr.HasError(err, MsgErrorCreatingUser)
}

func Login(cookieId string, username string, password string) error {
	username = strings.ToLower(username)
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return jerr.Get(MsgUsernameNotFound, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return jerr.Get(MsgPasswordMismatch, err)
	}

	session, err := db.GetSession(cookieId)
	if err != nil {
		return jerr.Get("session not found", err)
	}

	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return jerr.Get("session save failed", err)
	}

	return nil
}
