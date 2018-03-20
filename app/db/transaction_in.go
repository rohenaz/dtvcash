package db

import (
	"encoding/hex"
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"strings"
	"time"
)

type TransactionIn struct {
	Id                    uint   `gorm:"primary_key"`
	TransactionId         uint   `gorm:"unique_index:transaction_in_script;"`
	Transaction           *Transaction
	PreviousOutPointHash  []byte
	PreviousOutPointIndex uint32
	SignatureScript       []byte `gorm:"unique_index:transaction_in_script;"`
	UnlockString          string
	Witnesses             []*Witness
	Sequence              uint32
	TxnOutId              uint
	TxnOut                *TransactionOut
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (t TransactionIn) GetOutPoint() *wire.OutPoint {
	hash, _ := chainhash.NewHash(t.PreviousOutPointHash)
	return wire.NewOutPoint(hash, t.PreviousOutPointIndex)
}

func (t TransactionIn) GetPrevOutPointHash() []byte {
	return t.GetOutPoint().Hash.CloneBytes()
}

func (t TransactionIn) GetPrevOutPointString() string {
	return t.GetOutPoint().String()
}

func (t TransactionIn) HasOut() bool {
	return t.TxnOutId > 0
}

func (t TransactionIn) GetAddress() string {
	split := strings.Split(t.UnlockString, " ")
	if len(split) != 2 {
		return ""
	}
	pubKey, err := hex.DecodeString(split[1])
	if err != nil {
		return ""
	}
	return wallet.GetAddress(pubKey).GetEncoded()
}
