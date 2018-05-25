package db

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/rohenaz/dtvcash/app/bitcoin/wallet"
	"time"
)

type MemoPollVote struct {
	Id           uint   `gorm:"primary_key"`
	TxHash       []byte `gorm:"unique;size:50"`
	ParentHash   []byte
	PkHash       []byte `gorm:"index:pk_hash"`
	PkScript     []byte
	BlockId      uint
	Block        *Block
	OptionTxHash []byte `gorm:"index:option_tx_hash"`
	TipAmount    int64
	TipPkHash    []byte
	Message      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m MemoPollVote) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo poll vote", result.Error)
	}
	return nil
}

func (m MemoPollVote) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo poll vote", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPollVote) GetAddressString() string {
	return m.GetAddress().GetEncoded()
}

func (m MemoPollVote) GetAddress() wallet.Address {
	return wallet.GetAddressFromPkHash(m.PkHash)
}

func GetMemoPollVote(txHash []byte) (*MemoPollVote, error) {
	var memoPollVote MemoPollVote
	err := find(&memoPollVote, MemoPollVote{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo poll vote", err)
	}
	return &memoPollVote, nil
}

func GetVotesForOptions(questionTxHash []byte, single bool) ([]*MemoPollVote, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoPollVotes []*MemoPollVote
	if single {
		var joinSql = "JOIN (" +
			"SELECT MIN(memo_poll_votes.id) AS id " +
			"FROM memo_poll_votes " +
			"JOIN memo_poll_options ON (memo_poll_votes.option_tx_hash = memo_poll_options.tx_hash) " +
			"WHERE memo_poll_options.poll_tx_hash = ? " +
			"GROUP BY memo_poll_votes.pk_hash) AS uids ON (memo_poll_votes.id = uids.id)"
		db = db.Joins(joinSql, questionTxHash)
	} else {
		var joinSql = "JOIN (" +
			"SELECT memo_poll_votes.id AS id " +
			"FROM memo_poll_votes " +
			"JOIN memo_poll_options ON (memo_poll_votes.option_tx_hash = memo_poll_options.tx_hash) " +
			"WHERE memo_poll_options.poll_tx_hash = ? " +
			") AS uids ON (memo_poll_votes.id = uids.id)"
		db = db.Joins(joinSql, questionTxHash)
	}
	result := db.
		Order("created_at DESC").
		Find(&memoPollVotes)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo poll votes", result.Error)
	}
	return memoPollVotes, nil
}

func GetCountMemoPollVote() (uint, error) {
	cnt, err := count(&MemoPollVote{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
