package db

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"time"
)

type CsrfToken struct {
	Id         uint   `gorm:"primary_key"`
	CookieId   string `gorm:"unique;size:140"`
	Token      string `gorm:"unique;size:140"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c *CsrfToken) Save() error {
	result := save(&c)
	return result.Error
}

func GetCsrfTokenString(cookieId string) (string, error) {
	csrfToken, err := GetCsrfToken(cookieId)
	if err != nil {
		return "", jerr.Get("error getting csrf token", err)
	}
	return csrfToken.Token, nil
}

func GetCsrfToken(cookieId string) (*CsrfToken, error) {
	csrfToken := &CsrfToken{
		CookieId: cookieId,
	}
	err := find(csrfToken, csrfToken)
	if err == nil {
		return csrfToken, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting token from database", err)
	}
	csrfToken.Token = web.CreateToken()
	err = create(csrfToken)
	if err != nil {
		return nil, jerr.Get("error creating token", err)
	}
	return csrfToken, nil
}

func UpdateCsrfTokenSession(oldCookieId string, newCookieId string) error {
	csrfToken, err := GetCsrfToken(oldCookieId)
	if err != nil {
		return jerr.Get("error getting csrf token from db", err)
	}
	csrfToken.CookieId = newCookieId
	err = csrfToken.Save()
	if err != nil {
		return jerr.Get("error saving token", err)
	}
	return nil
}
