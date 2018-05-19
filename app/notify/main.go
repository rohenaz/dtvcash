package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type Notification interface {
	GetMessage() string
	GetLink() string
}

func GetNotificationsFeed(pkHash []byte, offset uint) ([]Notification, error) {
	dbNotifications, err := db.GetRecentNotificationsForUser(pkHash, offset)
	if err != nil {
		return nil, jerr.Get("error getting notifications from db", err)
	}
	var notifications []Notification
	for _, dbNotification := range dbNotifications {
		if dbNotification.IsTypeLike() {
			notifications = append(notifications, &LikeNotification{
				Notification: dbNotification,
			})
		}
	}
	return notifications, nil
}
