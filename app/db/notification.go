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
		Type:   notificationType,
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

func GetRecentNotificationsForUser(pkHash []byte, offset uint) ([]*Notification, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var notifications []*Notification
	result := db.
		Limit(25).
		Offset(offset).
		Joins("LEFT OUTER JOIN transactions ON (notifications.tx_hash = transactions.hash)").
		Joins("LEFT OUTER JOIN blocks ON (transactions.block_id = blocks.id)").
		Where("notifications.pk_hash = ?", pkHash).
		Order("COALESCE(blocks.timestamp, transactions.created_at) DESC, transactions.created_at DESC").
		Find(&notifications)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return notifications, nil
}

func GetUnreadNotificationCount(pkHash []byte, lastNotificationId uint) (uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	var count uint
	result := db.
		Table("notifications").
		Where("pk_hash = ?", pkHash).
		Where("id > ?", lastNotificationId).
		Count(&count)
	if result.Error != nil {
		return 0, jerr.Get("error getting unread notification count", result.Error)
	}
	return count, nil
}
