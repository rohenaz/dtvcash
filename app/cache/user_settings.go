package cache

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func GetUserSettings(userId uint) (*db.UserSettings, error) {
	var userSettings db.UserSettings
	err := GetItem(getUserSettingsName(userId), &userSettings)
	if err == nil {
		return &userSettings, nil
	}
	if ! IsMissError(err) {
		return nil, jerr.Get("error getting user settings from cache", err)
	}
	dbUserSettings, err := db.GetSettingsForUser(userId)
	if err != nil {
		return nil, jerr.Get("error getting user settings from db", err)
	}
	err = SetUserSettings(dbUserSettings)
	if err != nil {
		return nil, jerr.Get("error saving cache", err)
	}
	return dbUserSettings, nil
}

func SetUserSettings(settings *db.UserSettings) error {
	err := SetItem(getUserSettingsName(settings.UserId), settings)
	if err != nil {
		return jerr.Get("error setting user settings cache", err)
	}
	return nil
}

func getUserSettingsName(userId uint) string {
	return fmt.Sprintf("user-settings-%d", userId)
}
