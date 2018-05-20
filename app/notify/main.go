package notify

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type Notification struct {
	Generic Generic
}

func (n Notification) IsLike() bool {
	_, ok := n.Generic.(*LikeNotification)
	return ok
}

func (n Notification) IsReply() bool {
	_, ok := n.Generic.(*ReplyNotification)
	return ok
}

func (n Notification) GetPostMessage() string {
	var msg string
	switch g := n.Generic.(type) {
	case *LikeNotification:
		msg = g.Post.GetMessage()
	case *ReplyNotification:
		msg = g.Post.GetMessage()
	}
	if len(msg) > 50 {
		msg = msg[:47] + "..."
	}
	return msg
}

func (n Notification) GetParentMessage() string {
	var msg string
	switch g := n.Generic.(type) {
	case *ReplyNotification:
		msg = g.Parent.GetMessage()
	}
	if len(msg) > 50 {
		msg = msg[:47] + "..."
	}
	return msg
}

type Generic interface {
	GetMessage() string
	GetLink() string
	GetTime() string
}

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
				return nil, jerr.Get("error getting notification like", err)
			}
			post, err := db.GetMemoPost(like.LikeTxHash)
			if err != nil {
				return nil, jerr.Get("error getting like post", err)
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
				return nil, jerr.Get("error getting notification post", err)
			}
			parent, err := db.GetMemoPost(post.ParentTxHash)
			if err != nil {
				return nil, jerr.Get("error getting notification post parent", err)
			}
			generics = append(generics, &ReplyNotification{
				Notification: dbNotification,
				Post:         post,
				Parent:       parent,
			})
			namePkHashes = append(namePkHashes, post.PkHash)
		}
	}
	setNames, err := db.GetNamesForPkHashes(namePkHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hashes", err)
	}
	for _, notification := range generics {
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
	var notifications []*Notification
	for _, generic := range generics {
		notifications = append(notifications, &Notification{
			Generic: generic,
		})
	}
	return notifications, nil
}
