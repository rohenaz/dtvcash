package notify

import (
	"fmt"
	"time"
)

type Generic interface {
	GetName() string
	GetAddressString() string
	GetPostHashString() string
	GetMessage() string
	GetTime() time.Time
}

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

func (n Notification) IsNewFollower() bool {
	_, ok := n.Generic.(*NewFollowerNotification)
	return ok
}

func (n Notification) GetName() string {
	return n.Generic.GetName()
}

func (n Notification) GetAddressString() string {
	return n.Generic.GetAddressString()
}

func (n Notification) GetPostHashString() string {
	return n.Generic.GetPostHashString()
}

func (n Notification) GetParentHashString() string {
	reply, ok := n.Generic.(*ReplyNotification)
	if !ok {
		return ""
	}
	return reply.Parent.GetTransactionHashString()
}

func (n Notification) GetPostMessage() string {
	msg := n.Generic.GetMessage()
	return msg
}

func (n Notification) GetParentMessage() string {
	var msg string
	switch g := n.Generic.(type) {
	case *ReplyNotification:
		msg = g.Parent.GetMessage()
	}
	return msg
}

func (n Notification) GetTimeAgo() string {
	ts := n.Generic.GetTime()
	delta := time.Now().Sub(ts)
	hours := int(delta.Hours())
	if hours > 0 {
		if hours >= 24 {
			if hours < 48 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", hours/24)
		}
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	minutes := int(delta.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return fmt.Sprintf("%d seconds ago", int(delta.Seconds()))
}

func (n Notification) GetTipAmount() int64 {
	like, ok := n.Generic.(*LikeNotification)
	if !ok {
		return 0
	}
	return like.Like.TipAmount
}

func (n Notification) GetId() uint {
	switch g := n.Generic.(type) {
	case *ReplyNotification:
		return g.Notification.Id
	case *LikeNotification:
		return g.Notification.Id
	case *NewFollowerNotification:
		return g.Notification.Id
	}
	return 0
}
