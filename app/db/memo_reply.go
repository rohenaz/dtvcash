package db

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"time"
)

type MemoReply struct {
	Id          uint   `gorm:"primary_key"`
	TxHash      []byte `gorm:"unique;size:50"`
	ParentHash  []byte
	PkHash      []byte `gorm:"index:pk_hash"`
	PkScript    []byte
	Address     string
	ReplyTxHash []byte `gorm:"index:reply_tx_hash"`
	Message     string
	BlockId     uint
	Block       *Block
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m MemoReply) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoReply) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoReply) GetAddressString() string {
	return m.GetAddress().GetEncoded()
}

func (m MemoReply) GetAddress() wallet.Address {
	return wallet.GetAddressFromPkHash(m.PkHash)
}

func (m MemoReply) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoReply) GetMessage() string {
	return m.Message
}

func (m MemoReply) GetTimeString() string {
	if m.BlockId != 0 {
		if m.Block != nil {
			return m.Block.Timestamp.Format("2006-01-02 15:04:05")
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}

func GetMemoReply(txHash []byte) (*MemoReply, error) {
	var memoReply MemoReply
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoReply, MemoReply{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoReply, nil
}

type memoReplySortByDate []*MemoReply

func (txns memoReplySortByDate) Len() int      { return len(txns) }
func (txns memoReplySortByDate) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoReplySortByDate) Less(i, j int) bool {
	if bytes.Equal(txns[i].ParentHash, txns[j].TxHash) {
		return true
	}
	if bytes.Equal(txns[i].TxHash, txns[j].ParentHash) {
		return false
	}
	if txns[i].Block == nil && txns[j].Block == nil {
		return false
	}
	if txns[i].Block == nil {
		return true
	}
	if txns[j].Block == nil {
		return false
	}
	return txns[i].Block.Height > txns[j].Block.Height
}

func GetCountMemoReplies() (uint, error) {
	cnt, err := count(&MemoReply{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
