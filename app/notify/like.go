package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type LikeNotification struct {
	Notification *db.Notification
}

func (n LikeNotification) GetMessage() string {
	return ""
}

func (n LikeNotification) GetLink() string {
	return ""
}

func AddLikeNotification(like *db.MemoLike) error {
	post, err := db.GetMemoPost(like.LikeTxHash)
	if err != nil {
		return jerr.Get("error getting memo post", err)
	}
	_, err = db.AddNotification(post.PkHash, like.TxHash, db.NotificationTypeLike)
	if err != nil {
		return jerr.Get("error adding notification", err)
	}
	return nil
}
