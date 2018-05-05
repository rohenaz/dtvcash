package db

import (
	"bytes"
	"git.jasonc.me/main/memo/app/bitcoin/script"
	"git.jasonc.me/main/memo/app/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"sort"
	"time"
)

type MemoSetProfile struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte
	Address    string
	Profile    string
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoSetProfile) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoSetProfile) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoSetProfile) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoSetProfile) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoSetProfile) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoSetProfileById(id uint) (*MemoSetProfile, error) {
	var memoSetProfile MemoSetProfile
	err := find(&memoSetProfile, MemoSetProfile{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set profile", err)
	}
	return &memoSetProfile, nil
}

func GetMemoSetProfile(txHash []byte) (*MemoSetProfile, error) {
	var memoSetProfile MemoSetProfile
	err := find(&memoSetProfile, MemoSetProfile{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set profile", err)
	}
	return &memoSetProfile, nil
}

type memoSetProfileSortByDate []*MemoSetProfile

func (txns memoSetProfileSortByDate) Len() int                      { return len(txns) }
func (txns memoSetProfileSortByDate) Swap(i, j int)      { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoSetProfileSortByDate) Less(i, j int) bool {
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

func GetProfileForPkHash(pkHash []byte) (*MemoSetProfile, error) {
	profiles, err := GetSetProfilesForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting set profiles for pk hash", err)
	}
	if len(profiles) == 0 {
		return nil, nil
	}
	return profiles[0], nil
}

func GetSetProfilesForPkHash(pkHash []byte) ([]*MemoSetProfile, error) {
	var memoSetProfiles []*MemoSetProfile
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoSetProfiles, &MemoSetProfile{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo profiles", err)
	}
	sort.Sort(memoSetProfileSortByDate(memoSetProfiles))
	return memoSetProfiles, nil
}

func GetCountMemoSetProfile() (uint, error) {
	cnt, err := count(&MemoSetProfile{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
