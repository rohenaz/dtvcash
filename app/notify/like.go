package notify

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type LikeNotification struct {
	Like         *db.MemoLike
	Post         *db.MemoPost
	Notification *db.Notification
	Name         string
}

func (n LikeNotification) GetMessage() string {
	return fmt.Sprintf("%s liked your post", n.Name)
}

func (n LikeNotification) GetLink() string {
	hash, err := chainhash.NewHash(n.Like.LikeTxHash)
	if err != nil {
		jerr.Get("error getting like notification tx hash", err).Print()
		return ""
	}
	return fmt.Sprintf("post/%s", hash.String())
}

func (n LikeNotification) GetTime() string {
	return n.Like.GetTimeString()
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
