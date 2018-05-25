package notify

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func GetNotificationsFeed(pkHash []byte, offset uint) ([]*Notification, error) {
	dbNotifications, err := db.GetRecentNotificationsForUser(pkHash, offset)
	if err != nil {
		return nil, jerr.Get("error getting notifications from db", err)
	}
	var namePkHashes [][]byte
	var generics []Generic
	for _, dbNotification := range dbNotifications {
		switch dbNotification.Type {
		case db.NotificationTypeLike:
			like, err := db.GetMemoLike(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification like", err).Print()
				continue
			}
			post, err := db.GetMemoPost(like.LikeTxHash)
			if err != nil {
				jerr.Get("error getting like post for notification", err).Print()
				continue
			}
			generics = append(generics, &LikeNotification{
				Notification: dbNotification,
				Like:         like,
				Post:         post,
			})
			namePkHashes = append(namePkHashes, like.PkHash)
		case db.NotificationTypeReply:
			post, err := db.GetMemoPost(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification post", err).Print()
				continue
			}
			parent, err := db.GetMemoPost(post.ParentTxHash)
			if err != nil {
				jerr.Get("error getting notification post parent", err).Print()
				continue
			}
			generics = append(generics, &ReplyNotification{
				Notification: dbNotification,
				Post:         post,
				Parent:       parent,
			})
			namePkHashes = append(namePkHashes, post.PkHash)
		case db.NotificationTypeNewFollower:
			follow, err := db.GetMemoFollow(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification new follower", err).Print()
				continue
			}
			generics = append(generics, &NewFollowerNotification{
				Notification: dbNotification,
				Follow:       follow,
			})
			namePkHashes = append(namePkHashes, follow.PkHash)
		}
	}
	setNames, err := db.GetNamesForPkHashes(namePkHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hashes", err)
	}
	for _, notification := range generics {
		switch n := notification.(type) {
		case *LikeNotification:
			for _, setName := range setNames {
				if bytes.Equal(n.Like.PkHash, setName.PkHash) {
					n.Name = setName.Name
				}
			}
			if n.Name == "" {
				n.Name = n.Like.GetAddressString()[:16]
			}
		case *ReplyNotification:
			for _, setName := range setNames {
				if bytes.Equal(n.Post.PkHash, setName.PkHash) {
					n.Name = setName.Name
				}
			}
			if n.Name == "" {
				n.Name = n.Post.GetAddressString()[:16]
			}
		case *NewFollowerNotification:
			for _, setName := range setNames {
				if bytes.Equal(n.Follow.PkHash, setName.PkHash) {
					n.Name = setName.Name
				}
			}
			if n.Name == "" {
				n.Name = n.Follow.GetAddressString()[:16]
			}
		}
	}
	var notifications []*Notification
	for _, generic := range generics {
		notifications = append(notifications, &Notification{
			Generic: generic,
		})
	}
	return notifications, nil
}
