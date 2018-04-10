package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type User struct {
	Id           uint   `gorm:"primary_key"`
	Username     string `gorm:"unique;size:50"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateUser(username string, hashedPassword string) (*User, error) {
	user := &User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	err := create(user)
	if err != nil {
		return nil, jerr.Get("error creating user", err)
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{
		Username: username,
	}
	err := find(user, user)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func GetUserById(id uint) (*User, error) {
	user := &User{
		Id: id,
	}
	err := find(user, user)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}
