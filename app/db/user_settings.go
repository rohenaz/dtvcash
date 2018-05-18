package db

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const (
	SettingIntegrationsAll  = "all"
	SettingIntegrationsHide = "hide"
	SettingIntegrationsNone = "none"

	SettingThemeDefault = "default"
	SettingThemeDark    = "dark"

	MaxDefaultTip = 1e9
)

type UserSettings struct {
	Id           uint   `gorm:"primary_key"`
	UserId       uint   `gorm:"unique"`
	DefaultTip   uint
	Integrations string `gorm:"size:25"`
	Theme        string `gorm:"size:25"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *UserSettings) Save() error {
	result := save(u)
	if result.Error != nil {
		return jerr.Get("error saving user settings", result.Error)
	}
	return nil
}

func (u UserSettings) GetDefaultTipString() string {
	if u.DefaultTip == 0 {
		return ""
	}
	return fmt.Sprintf("%d", u.DefaultTip)
}

func SaveSettingsForUser(userId uint, defaultTip uint, integrations string, theme string) (*UserSettings, error) {
	var userSettings = UserSettings{
		UserId: userId,
	}
	err := find(&userSettings, userSettings)
	if err != nil && ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error finding existing settings", err)
	}
	userSettings.DefaultTip = defaultTip
	userSettings.Integrations = integrations
	userSettings.Theme = theme
	err = userSettings.Save()
	if err != nil {
		return nil, jerr.Get("error saving settings", err)
	}
	return &userSettings, nil
}

func GetSettingsForUser(userId uint) (*UserSettings, error) {
	var userSettings = UserSettings{
		UserId: userId,
	}
	err := find(&userSettings, userSettings)
	if err == nil {
		return &userSettings, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error finding settings for user", err)
	}
	// Defaults
	userSettings.Integrations = SettingIntegrationsAll
	userSettings.Theme = SettingThemeDefault
	err = userSettings.Save()
	if err != nil {
		return nil, jerr.Get("error saving default settings", err)
	}
	return &userSettings, nil
}

func IsValidDefaultTip(defaultTip uint) bool {
	return defaultTip == 0 || (defaultTip >= 546 && defaultTip <= MaxDefaultTip)
}

func IsValidIntegrationsSetting(integrations string) bool {
	for _, validValue := range []string{
		SettingIntegrationsAll,
		SettingIntegrationsHide,
		SettingIntegrationsNone,
	} {
		if integrations == validValue {
			return true
		}
	}
	return false
}

func IsValidThemeSetting(theme string) bool {
	for _, validValue := range []string{
		SettingThemeDefault,
		SettingThemeDark,
	} {
		if theme == validValue {
			return true
		}
	}
	return false
}
