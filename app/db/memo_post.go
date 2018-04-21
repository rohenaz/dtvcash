package db

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"sort"
	"time"
)

type MemoPost struct {
	Id           uint        `gorm:"primary_key"`
	TxHash       []byte      `gorm:"unique;size:50"`
	ParentHash   []byte
	PkHash       []byte      `gorm:"index:pk_hash"`
	PkScript     []byte
	Address      string
	ParentTxHash []byte      `gorm:"index:parent_tx_hash"`
	Parent       *MemoPost
	Replies      []*MemoPost `gorm:"foreignkey:ParentTxHash"`
	Message      string
	BlockId      uint
	Block        *Block
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	return m.GetAddress().GetEncoded()
}

func (m MemoPost) GetAddress() wallet.Address {
	return wallet.GetAddressFromPkHash(m.PkHash)
}

func (m MemoPost) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoPost) GetMessage() string {
	return m.Message
}

func (m MemoPost) GetTimeString() string {
	if m.BlockId != 0 {
		if m.Block != nil {
			return m.Block.Timestamp.Format("2006-01-02 15:04:05")
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}

func GetMemoPost(txHash []byte) (*MemoPost, error) {
	var memoPost MemoPost
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoPost, MemoPost{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoPost, nil
}

type memoPostSortByDate []*MemoPost

func (txns memoPostSortByDate) Len() int      { return len(txns) }
func (txns memoPostSortByDate) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoPostSortByDate) Less(i, j int) bool {
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

func GetPostsForPkHashes(pkHashes [][]byte) ([]*MemoPost, error) {
	if len(pkHashes) == 0 {
		return nil, nil
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.Preload(BlockTable).Where("pk_hash in (?)", pkHashes).Where("block_id > 0").Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetPostsForPkHash(pkHash []byte) ([]*MemoPost, error) {
	if len(pkHash) == 0 {
		return nil, nil
	}
	var memoPosts []*MemoPost
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoPosts, &MemoPost{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo posts", err)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetUniqueMemoAPkHashes() ([][]byte, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	rows, err := db.Table("memo_posts").Select("DISTINCT(pk_hash)").Rows()
	if err != nil {
		return nil, jerr.Get("error getting distinct pk hashes", err)
	}
	var pkHashes [][]byte
	for rows.Next() {
		var pkHash []byte
		err := rows.Scan(&pkHash)
		if err != nil {
			return nil, jerr.Get("error scanning row with pkHash", err)
		}
		pkHashes = append(pkHashes, pkHash)
	}
	return pkHashes, nil
}

func GetRecentPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.Limit(25).Offset(offset).Order("id DESC").Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetCountMemoPosts() (uint, error) {
	cnt, err := count(&MemoPost{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
