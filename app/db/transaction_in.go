package db

import (
	"encoding/hex"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"strings"
	"time"
)

type TransactionIn struct {
	Id                    uint   `gorm:"primary_key"`
	Index                 uint
	TransactionId         uint   `gorm:"unique_index:transaction_in_script;"`
	Transaction           *Transaction
	KeyId                 uint
	Key                   *Key
	PreviousOutPointHash  []byte
	PreviousOutPointIndex uint32
	SignatureScript       []byte `gorm:"unique_index:transaction_in_script;"`
	UnlockString          string
	Sequence              uint32
	TxnOutId              uint
	TxnOut                *TransactionOut
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (t TransactionIn) Save() error {
	result := save(&t)
	if result.Error != nil {
		return jerr.Get("error saving transaction input", result.Error)
	}
	return nil
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

func (t TransactionIn) GetAddressString() string {
	return t.GetAddress().GetEncoded()
}
func (t TransactionIn) GetAddress() wallet.Address {
	split := strings.Split(t.UnlockString, " ")
	if len(split) != 2 {
		return wallet.Address{}
	}
	pubKey, err := hex.DecodeString(split[1])
	if err != nil {
		return wallet.Address{}
	}
	return wallet.GetAddress(pubKey)
}

func (t TransactionIn) Delete() error {
	if t.TxnOutId != 0 {
		var txOut TransactionOut
		err := find(&txOut, TransactionOut{Id: t.TxnOutId})
		if err != nil {
			if ! IsRecordNotFoundError(err) {
				return jerr.Get("error getting transaction out", err)
			}
		} else {
			txOut.TxnInId = 0
			err = txOut.Save()
			if err != nil {
				return jerr.Get("error saving transaction out", err)
			}
		}
	}
	result := remove(t)
	if result.Error != nil {
		return jerr.Get("error removing transaction input", result.Error)
	}
	return nil
}

func GetTransactionInputById(id uint) (*TransactionIn, error) {
	query, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var txIn TransactionIn
	result := query.
		Preload("Transaction").
		Preload("Transaction.Block").
		Find(&txIn, TransactionIn{
		Id: id,
	})
	if result.Error != nil {
		return nil, jerr.Get("error finding transaction in", result.Error)
	}
	return &txIn, nil
}
