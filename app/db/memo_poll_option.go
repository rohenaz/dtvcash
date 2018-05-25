package db

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"time"
)

type MemoPollOption struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte
	PollTxHash []byte `gorm:"index:poll_tx_hash"`
	Option     string `gorm:"size:500"`
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoPollOption) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo poll option", result.Error)
	}
	return nil
}

func (m MemoPollOption) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo poll option", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPollOption) GetPollTransactionHashString() string {
	hash, err := chainhash.NewHash(m.PollTxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo poll option", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPollOption) GetAddressString() string {
	return m.GetAddress().GetEncoded()
}

func (m MemoPollOption) GetAddress() wallet.Address {
	return wallet.GetAddressFromPkHash(m.PkHash)
}

func GetMemoPollOption(txHash []byte) (*MemoPollOption, error) {
	var memoPollOption MemoPollOption
	err := find(&memoPollOption, MemoPollOption{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo poll option", err)
	}
	return &memoPollOption, nil
}

func GetMemoPollOptionByOption(pollTxHash []byte, option string) (*MemoPollOption, error) {
	var memoPollOption MemoPollOption
	err := find(&memoPollOption, MemoPollOption{
		PollTxHash: pollTxHash,
		Option:     option,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo poll option", err)
	}
	return &memoPollOption, nil
}
