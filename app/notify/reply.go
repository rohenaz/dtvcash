package notify

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type ReplyNotification struct {
	Post         *db.MemoPost
	Parent       *db.MemoPost
	Notification *db.Notification
	Name         string
}

func (n ReplyNotification) GetMessage() string {
	return fmt.Sprintf("%s replied to your post", n.Name)
}

func (n ReplyNotification) GetLink() string {
	hash, err := chainhash.NewHash(n.Post.TxHash)
	if err != nil {
		jerr.Get("error getting like notification tx hash", err).Print()
		return ""
	}
	return fmt.Sprintf("post/%s", hash.String())
}

func (n ReplyNotification) GetTime() string {
	if n.Post.BlockId != 0 {
		if n.Post.Block != nil {
			return n.Post.Block.Timestamp.Format("2006-01-02 15:04:05")
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}

func AddReplyNotification(reply *db.MemoPost) error {
	parent, err := db.GetMemoPost(reply.ParentTxHash)
	if err != nil {
		return jerr.Get("error getting parent post", err)
	}
	_, err = db.AddNotification(parent.PkHash, reply.TxHash, db.NotificationTypeReply)
	if err != nil {
		return jerr.Get("error adding notification", err)
	}
	return nil
}
