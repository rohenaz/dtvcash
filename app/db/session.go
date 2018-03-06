package db

import (
	"strings"
	"time"
)

type Session struct {
	Id           uint   `gorm:"primary_key"`
	CookieId     string `gorm:"unique_index"`
	HasLoggedOut bool
	UserId       uint
	StartTs      uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *Session) Save() error {
	result := save(&s)
	return result.Error
}

func GetSession(cookieId string) (*Session, error) {
	session := &Session{
		CookieId: cookieId,
	}
	err := find(session, session)
	if err != nil && strings.Contains(err.Error(), "record not found") {
		err = create(session)
	}
	if err != nil {
		return nil, err
	} else {
		return session, nil
	}
}
