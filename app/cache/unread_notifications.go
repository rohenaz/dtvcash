package cache

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type UnreadNotifications struct {
	Count uint
}

func GetUnreadNotificationCount(userId uint) (uint, error) {
	var unreadNotifications UnreadNotifications
	err := GetItem(getUnreadNotificationName(userId), &unreadNotifications)
	if err == nil {
		return unreadNotifications.Count, nil
	}
	if ! IsMissError(err) {
		return 0, jerr.Get("error getting pk hash from cache", err)
	}
	unreadCount, err := GetAndSetUnreadNotificationCount(userId)
	if err != nil {
		return 0, jerr.Get("error setting user last notification id cache", err)
	}
	return unreadCount, nil
}

func GetAndSetUnreadNotificationCount(userId uint) (uint, error) {
	lastNotificationId, err := db.GetLastNotificationId(userId)
	if err != nil {
		return 0, jerr.Get("error getting last notification id from db", err)
	}
	pkHash, err := GetUserPkHash(userId)
	if err != nil {
		return 0, jerr.Get("error getting user pk hash", err)
	}
	unreadCount, err := db.GetUnreadNotificationCount(pkHash, lastNotificationId)
	if err != nil {
		return 0, jerr.Get("error getting unread count", err)
	}
	err = SetUnreadNotificationCount(userId, unreadCount)
	if err != nil {
		return 0, jerr.Get("error setting unread notification count in cache", err)
	}
	return unreadCount, nil
}

func SetUnreadNotificationCount(userId uint, count uint) error {
	err := SetItem(getUnreadNotificationName(userId), UnreadNotifications{
		Count: count,
	})
	if err != nil {
		return jerr.Get("error setting user last notification id cache", err)
	}
	return nil
}

func getUnreadNotificationName(userId uint) string {
	return fmt.Sprintf("user-unread-notifications-%d", userId)
}
