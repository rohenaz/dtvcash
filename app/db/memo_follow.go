package db

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"sort"
	"time"
)

type MemoFollow struct {
	Id           uint   `gorm:"primary_key"`
	TxHash       []byte `gorm:"unique;size:50"`
	ParentHash   []byte
	PkHash       []byte `gorm:"index:pk_hash"`
	PkScript     []byte
	Address      string
	FollowPkHash []byte `gorm:"index:follow_pk_hash"`
	BlockId      uint
	Block        *Block
	Unfollow     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m MemoFollow) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoFollow) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoFollow) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoFollow) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoFollow) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoFollow(txHash []byte) (*MemoFollow, error) {
	var memoFollow MemoFollow
	err := find(&memoFollow, MemoFollow{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoFollow, nil
}

type memoFollowSortByDate []*MemoFollow

func (txns memoFollowSortByDate) Len() int      { return len(txns) }
func (txns memoFollowSortByDate) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoFollowSortByDate) Less(i, j int) bool {
	if bytes.Equal(txns[i].ParentHash, txns[j].TxHash) {
		return false
	}
	if bytes.Equal(txns[i].TxHash, txns[j].ParentHash) {
		return true
	}
	if txns[i].Block == nil && txns[j].Block == nil {
		return false
	}
	if txns[i].Block == nil {
		return false
	}
	if txns[j].Block == nil {
		return true
	}
	return txns[i].Block.Height < txns[j].Block.Height
}

func GetFollowersForPkHash(pkHash []byte) ([]*MemoFollow, error) {
	var memoFollows []*MemoFollow
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoFollows, &MemoFollow{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo follows", err)
	}
	sort.Sort(memoFollowSortByDate(memoFollows))
	return memoFollows, nil
}

func GetFollowingForPkHash(followPkHash []byte) ([]*MemoFollow, error) {
	var memoFollows []*MemoFollow
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoFollows, &MemoFollow{
		FollowPkHash: followPkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo follows", err)
	}
	sort.Sort(memoFollowSortByDate(memoFollows))
	return memoFollows, nil
}

func GetFollowingCountForPkHash(pkHash []byte) (uint, error) {
	cnt, err := count(&MemoFollow{
		PkHash: pkHash,
	})
	if err != nil {
		return 0, jerr.Get("error getting follower count", err)
	}
	return cnt, nil
}

func GetFollowerCountForPkHash(followPkHash []byte) (uint, error) {
	cnt, err := count(&MemoFollow{
		FollowPkHash: followPkHash,
	})
	if err != nil {
		return 0, jerr.Get("error getting follower count", err)
	}
	return cnt, nil
}

func GetCountMemoFollows() (uint, error) {
	cnt, err := count(&MemoFollow{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
