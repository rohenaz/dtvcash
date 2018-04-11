package db

import (
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"time"
)

type MemoPost struct {
	Id        uint   `gorm:"primary_key"`
	TxHash    []byte `gorm:"unique;size:50"`
	PkHash    []byte
	PkScript  []byte
	Address   string
	Message   string
	BlockId   uint
	Block     *Block
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m MemoPost) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoPost) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPost) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoPost) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoPost) GetMessage() string {
	return html.EscapeString(m.Message)
}

func (m MemoPost) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoPost(txHash []byte) (*MemoPost, error) {
	var memoPost MemoPost
	err := find(&memoPost, MemoPost{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoPost, nil
}

func GetPostsForPkHash(pkHash []byte) ([]*MemoPost, error) {
	var memoPosts []*MemoPost
	err := findPreloadColumns([]string{
		BlockTable,
	},&memoPosts, &MemoPost{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo posts", err)
	}
	return memoPosts, nil
}
