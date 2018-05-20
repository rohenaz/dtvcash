package notify

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type Notification interface {
	GetMessage() string
	GetLink() string
	GetTime() string
}

func (n Notification) IsLike() bool {

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
				Like:         like,
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
	setNames, err := db.GetNamesForPkHashes(namePkHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hashes", err)
	}
	for _, notification := range notifications {
		for _, setName := range setNames {
			switch n := notification.(type) {
			case *LikeNotification:
				if bytes.Equal(n.Like.PkHash, setName.PkHash) {
					fmt.Printf("-")
					n.Name = setName.Name
				}
			case *ReplyNotification:
				if bytes.Equal(n.Post.PkHash, setName.PkHash) {
					fmt.Printf("_")
					n.Name = setName.Name
				}
			}
		}
	}

	return notifications, nil
}
