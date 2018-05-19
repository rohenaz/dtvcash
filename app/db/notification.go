package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const (
	NotificationTypeLike        = 1
	NotificationTypeReply       = 2
	NotificationTypeThreadReply = 3
)

type Notification struct {
	Id        uint   `gorm:"primary_key"`
	PkHash    []byte `gorm:"not null;unique_index:pk_hash_tx_hash"`
	TxHash    []byte `gorm:"not null;unique_index:pk_hash_tx_hash"`
	Type      uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AddNotification(pkHash []byte, txHash []byte, notificationType uint) (*Notification, error) {
	var notification = Notification{
		PkHash: pkHash,
		TxHash: txHash,
		Type: notificationType,
	}
	err := create(&notification)
	if err == nil {
		return &notification, nil
	}
	if ! IsDuplicateEntryError(err) {
		return nil, jerr.Get("error creating db notification", err)
	}
	return nil, nil
}
