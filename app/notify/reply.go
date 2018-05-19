package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

type ReplyNotification struct {
	Post         *db.MemoPost
	Notification *db.Notification
}

func (n ReplyNotification) GetMessage() string {
	return ""
}

func (n ReplyNotification) GetLink() string {
	return ""
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
