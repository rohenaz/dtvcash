package db

import (
	"git.jasonc.me/main/bitcoin/app/db"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const (
	StartScanBlock = 525000
)

type NodeStatus struct {
	Id        uint `gorm:"primary_key"`
	LastBlock uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *NodeStatus) Save() error {
	result := save(&s)
	return result.Error
}

func GetNodeStatus() (*NodeStatus, error) {
	status := &NodeStatus{
		Id: 1,
	}
	err := find(status, status)
	if err == nil {
		return status, nil
	}
	if ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting status", err)
	}
	status.LastBlock = StartScanBlock
	err = create(status)
	if err != nil {
		return nil, jerr.Get("error creating status", err)
	}
	return status, nil
}
