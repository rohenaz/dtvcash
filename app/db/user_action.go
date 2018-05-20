package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type UserAction struct {
	Id                 uint `gorm:"primary_key"`
	UserId             uint
	LastNotificationId uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func GetLastNotificationId(userId uint) (uint, error) {
	var userAction UserAction
	err := find(&userAction, UserAction{
		UserId: userId,
	})
	if err == nil {
		return userAction.LastNotificationId, nil
	}
	if ! IsRecordNotFoundError(err) {
		return 0, jerr.Get("error finding last notification", err)
	}
	return 0, nil
}

func SetLastNotificationId(userId uint, lastNotificationId uint) error {
	var userAction UserAction
	err := find(&userAction, UserAction{
		UserId: userId,
	})
	if err != nil {
		if ! IsRecordNotFoundError(err) {
			return jerr.Get("error getting last user action from db", err)
		}
		userAction = UserAction{
			UserId:             userId,
			LastNotificationId: lastNotificationId,
		}
		err := create(&userAction)
		if err != nil {
			return jerr.Get("error creating user action", err)
		}
		return nil
	} else {
		userAction.LastNotificationId = lastNotificationId
		result := save(userAction)
		if result.Error != nil {
			return jerr.Get("error saving user action", err)
		}
		return nil
	}
}
