package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Settings struct {
	Id            uint `gorm:"primary_key"`
	UserId        uint
	DefaultFilter uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s *Settings) Save() error {
	if s.Id == 0 {
		return jerr.New("must be an existing settings")
	}
	result := save(&s)
	if result.Error != nil {
		return jerr.Get("error saving settings", result.Error)
	}
	return nil
}

func (s *Settings) IsDefaultFilter(filterId uint) bool {
	return s.DefaultFilter == filterId
}

func GetSettings(userId uint) (*Settings, error) {
	settings := &Settings{
		UserId: userId,
	}
	err := find(settings, settings)
	if err == nil {
		return settings, nil
	}
	if ! isRecordNotFoundError(err) {
		return nil, jerr.Get("error looking up settings in database", err)
	}
	err = create(settings)
	if err != nil {
		return nil, jerr.Get("error creating settings", err)
	}
	return settings, nil
}
