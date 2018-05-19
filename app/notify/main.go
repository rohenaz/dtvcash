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
	var namePkHashes [][]byte
	var notifications []Notification
	for _, dbNotification := range dbNotifications {
		switch dbNotification.Type {
		case db.NotificationTypeLike:
			like, err := db.GetMemoLike(dbNotification.TxHash)
			if err != nil {
				return nil, jerr.Get("error getting notification like", err)
			}
			notifications = append(notifications, &LikeNotification{
				Notification: dbNotification,
				MemoLike:     like,
			})
			namePkHashes = append(namePkHashes, like.PkHash)
		case db.NotificationTypeReply:
			post, err := db.GetMemoPost(dbNotification.TxHash)
			if err != nil {
				return nil, jerr.Get("error getting notification post", err)
			}
			notifications = append(notifications, &ReplyNotification{
				Notification: dbNotification,
				Post: post,
			})
			namePkHashes = append(namePkHashes, post.PkHash)
		}
	}
	//setNames, err := db.GetNamesForPkHashes(namePkHashes)

	return notifications, nil
}
