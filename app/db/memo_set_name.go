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

type MemoSetName struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte
	Address    string
	Name    string
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoSetName) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoSetName) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoSetName) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoSetName) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoSetName) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoSetName(txHash []byte) (*MemoSetName, error) {
	var memoSetName MemoSetName
	err := find(&memoSetName, MemoSetName{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoSetName, nil
}

type memoSetNameSortByDate []*MemoSetName

func (txns memoSetNameSortByDate) Len() int                      { return len(txns) }
func (txns memoSetNameSortByDate) Swap(i, j int)      { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoSetNameSortByDate) Less(i, j int) bool {
	if bytes.Equal(txns[i].ParentHash, txns[j].TxHash) {
		return true
	}
	if bytes.Equal(txns[i].TxHash, txns[j].ParentHash) {
		return false
	}
	if txns[i].Block == nil && txns[j].Block == nil{
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

func GetNameForPkHash(pkHash []byte) (*MemoSetName, error) {
	names, err := GetSetNamesForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hash", err)
	}
	if len(names) == 0 {
		return nil, nil
	}
	return names[0], nil
}

func GetSetNamesForPkHash(pkHash []byte) ([]*MemoSetName, error) {
	var memoSetNames []*MemoSetName
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoSetNames, &MemoSetName{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo names", err)
	}
	sort.Sort(memoSetNameSortByDate(memoSetNames))
	return memoSetNames, nil
}

func GetCountMemoSetName() (uint, error) {
	cnt, err := count(&MemoSetName{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
