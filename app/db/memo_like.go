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

type MemoLike struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte
	Address    string
	LikeTxHash []byte `gorm:"index:like_tx_hash"`
	TipAmount  int64
	TipPkHash  []byte
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoLike) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoLike) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo like tx_hash", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoLike) GetLikeTransactionHashString() string {
	hash, err := chainhash.NewHash(m.LikeTxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo like like_tx_hash", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoLike) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoLike) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoLike) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoLike(txHash []byte) (*MemoLike, error) {
	var memoLike MemoLike
	err := find(&memoLike, MemoLike{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoLike, nil
}

type memoLikeSortByDate []*MemoLike

func (txns memoLikeSortByDate) Len() int      { return len(txns) }
func (txns memoLikeSortByDate) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoLikeSortByDate) Less(i, j int) bool {
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
		return true
	}
	if txns[j].Block == nil {
		return false
	}
	return txns[i].Block.Height < txns[j].Block.Height
}

func GetMemoLikesForTxnHash(txHash []byte) ([]*MemoLike, error) {
	var memoLikes []*MemoLike
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoLikes, &MemoLike{
		LikeTxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo likes", err)
	}
	sortReverse(memoLikeSortByDate(memoLikes))
	return memoLikes, nil
}

func GetMemoLikesForPkHash(pkHash []byte) ([]*MemoLike, error) {
	if len(pkHash) == 0 {
		return nil, nil
	}
	var memoLikes []*MemoLike
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoLikes, &MemoLike{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo likes", err)
	}
	sortReverse(memoLikeSortByDate(memoLikes))
	return memoLikes, nil
}

func sortReverse(memoLikes []*MemoLike) {
	sort.Sort(memoLikeSortByDate(memoLikes))
	for i, j := 0, len(memoLikes)-1; i < j; i, j = i+1, j-1 {
		memoLikes[i], memoLikes[j] = memoLikes[j], memoLikes[i]
	}
}

func GetCountMemoLikes() (uint, error) {
	cnt, err := count(&MemoLike{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}

func GetRecentTopLikedTxHashes(offset uint, timeStart time.Time, timeEnd time.Time) ([][]byte, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_likes").
		Select("like_tx_hash, COUNT(DISTINCT pk_hash) AS count").
		Joins("LEFT OUTER JOIN blocks ON (memo_likes.block_id = blocks.id)").
		Group("like_tx_hash").
		Order("count DESC, memo_likes.id DESC").
		Limit(25).
		Offset(offset)
	if timeEnd.IsZero() {
		query = query.Where("timestamp >= ? OR timestamp IS NULL", timeStart)
	} else {
		query = query.Where("timestamp >= ?", timeStart).Where("timestamp < ?", timeEnd)
	}
	rows, err := query.Rows()
	if err != nil {
		return nil, jerr.Get("error running query", err)
	}
	defer rows.Close()
	var txHashes [][]byte
	for rows.Next() {
		var likeTxHash []byte
		var count uint
		err := rows.Scan(&likeTxHash, &count)
		if err != nil {
			return nil, jerr.Get("error scanning rows", err)
		}
		txHashes = append(txHashes, likeTxHash)
	}
	return txHashes, nil
}

func GetPersonalizedRecentTopLikedTxHashes(selfPkHash []byte, offset uint, timeStart time.Time, timeEnd time.Time) ([][]byte, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "SELECT " +
		"	follow_pk_hash " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id) " +
		"WHERE unfollow = 0"
	query := db.
		Table("memo_likes").
		Select("like_tx_hash, COUNT(DISTINCT pk_hash) AS count").
		Joins("LEFT OUTER JOIN blocks ON (memo_likes.block_id = blocks.id)").
		Joins("JOIN (" + joinSelect + ") fsq ON (memo_likes.pk_hash = fsq.follow_pk_hash)", selfPkHash).
		Group("like_tx_hash").
		Order("count DESC, memo_likes.id DESC").
		Limit(25).
		Offset(offset)
	if timeEnd.IsZero() {
		query = query.Where("timestamp >= ? OR timestamp IS NULL", timeStart)
	} else {
		query = query.Where("timestamp >= ?", timeStart).Where("timestamp < ?", timeEnd)
	}
	rows, err := query.Rows()
	if err != nil {
		return nil, jerr.Get("error running query", err)
	}
	defer rows.Close()
	var txHashes [][]byte
	for rows.Next() {
		var likeTxHash []byte
		var count uint
		err := rows.Scan(&likeTxHash, &count)
		if err != nil {
			return nil, jerr.Get("error scanning rows", err)
		}
		txHashes = append(txHashes, likeTxHash)
	}
	return txHashes, nil
}

func GetRecentLikesForTopic(topic string, lastLikeId uint) ([]*MemoLike, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoLikes []*MemoLike
	result := db.
		Table("memo_likes").
		Select("memo_likes.*").
		Joins("JOIN memo_posts ON (memo_likes.like_tx_hash = memo_posts.tx_hash)").
		Where("memo_likes.id > ?", lastLikeId).
		Where("memo_posts.topic = ?", topic).
		Order("id ASC").
		Find(&memoLikes)
	if result.Error != nil {
		return nil, jerr.Get("error running recent topic likes query", result.Error)
	}
	return memoLikes, nil
}
