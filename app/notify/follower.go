package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"time"
)

type NewFollowerNotification struct {
	Follow       *db.MemoFollow
	Notification *db.Notification
	Name         string
}

func (n NewFollowerNotification) GetName() string {
	return n.Name
}

func (n NewFollowerNotification) GetPostHashString() string {
	return ""
}

func (n NewFollowerNotification) GetAddressString() string {
	return n.Follow.GetAddressString()
}

func (n NewFollowerNotification) GetMessage() string {
	return ""
}

func (n NewFollowerNotification) GetTime() time.Time {
	if n.Follow.Block != nil && n.Follow.Block.Timestamp.Before(n.Follow.CreatedAt) {
		return n.Follow.Block.Timestamp
	} else {
		return n.Follow.CreatedAt
	}
}

func AddNewFollowerNotification(follow *db.MemoFollow, updateCache bool) error {
	userId, err := db.GetUserIdFromPkHash(follow.FollowPkHash)
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			// Don't add notifications for external users, not an error though
			return nil
		}
		return jerr.Get("error getting user id from pk hash", err)
	}
	_, err = db.AddNotification(follow.FollowPkHash, follow.TxHash, db.NotificationTypeNewFollower)
	if err != nil {
		return jerr.Get("error adding notification", err)
	}
	if updateCache {
		_, err = cache.GetAndSetUnreadNotificationCount(userId)
		if err != nil {
			return jerr.Get("error setting notification unread count", err)
		}
	}
	return nil
}
