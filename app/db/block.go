package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Block struct {
	Id         uint `gorm:"primary_key"`
	Height     uint
	Timestamp  time.Time
	Hash       []byte
	PrevBlock  []byte
	MerkleRoot []byte
	Nonce      uint
	TxnCount   uint
	Version    int32
	Difficulty uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func GetGenesis() (*Block, error) {
	block := Block{
		Hash:       wallet.GenesisBlock.Hash.CloneBytes(),
		Timestamp:  time.Unix(1231006505, 0),
		MerkleRoot: wallet.GenesisBlock.MerkleRoot.CloneBytes(),
	}
	err := find(&block, &block)
	if err == nil {
		return &block, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error finding genesis block", err)
	}
	err = create(&block)
	if err != nil {
		return nil, jerr.Get("error creating block", err)
	}
	return &block, nil
}
