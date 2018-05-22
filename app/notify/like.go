package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"time"
)

type LikeNotification struct {
	Like         *db.MemoLike
	Post         *db.MemoPost
	Notification *db.Notification
	Name         string
}

func (n LikeNotification) GetName() string {
	return n.Name
}

func (n LikeNotification) GetAddressString() string {
	return n.Like.GetAddressString()
}

func (n LikeNotification) GetPostHashString() string {
	return n.Like.GetLikeTransactionHashString()
}

func (n LikeNotification) GetMessage() string {
	return n.Post.GetMessage()
}

func (n LikeNotification) GetTime() time.Time {
	if n.Like.Block != nil {
		return n.Like.Block.Timestamp
	} else {
		return n.Like.CreatedAt
	}
}

func AddLikeNotification(like *db.MemoLike, updateCache bool) error {
	post, err := db.GetMemoPost(like.LikeTxHash)
	if err != nil {
		return jerr.Get("error getting memo post", err)
	}
	userId, err := db.GetUserIdFromPkHash(post.PkHash)
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			// Don't add notifications for external users, not an error though
			return nil
		}
		return jerr.Get("error getting user id from pk hash", err)
	}
	_, err = db.AddNotification(post.PkHash, like.TxHash, db.NotificationTypeLike)
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