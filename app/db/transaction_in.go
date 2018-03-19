package db

import (
	"encoding/hex"
	"strconv"
	"time"
)

type TransactionIn struct {
	Id                    uint   `gorm:"primary_key"`
	TransactionId         uint
	PreviousOutPointHash  []byte
	PreviousOutPointIndex uint32
	SignatureScript       []byte `gorm:"unique;"`
	UnlockString          string
	Witnesses             []*Witness
	Sequence              uint32
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (t TransactionIn) GetPrevOutPointString() string {
	return hex.EncodeToString(t.PreviousOutPointHash) + ":" + strconv.Itoa(int(t.PreviousOutPointIndex))
}
