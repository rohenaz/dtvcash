package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

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
