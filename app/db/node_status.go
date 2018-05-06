package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type NodeStatus struct {
	Id            uint `gorm:"primary_key"`
	HeightChecked uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting status", err)
	}
	err = create(status)
	if err != nil {
		return nil, jerr.Get("error creating status", err)
	}
	return status, nil
}
