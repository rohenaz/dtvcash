package db

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/gorm"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"html"
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
		jerr.Get("error getting chainhash from memo follow", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoFollow) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo follow", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoFollow) GetFollowAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.FollowPkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo follow", err).Print()
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

func GetFollowersForPkHash(pkHash []byte, offset int) ([]*MemoFollow, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	sql := "" +
		"SELECT " +
		"	memo_follows.* " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id) " +
		"WHERE unfollow = 0 "
	var query *gorm.DB
	if offset >= 0 {
		sql += "LIMIT ?,25"
		query = db.Raw(sql, pkHash, offset)
	} else {
		query = db.Raw(sql, pkHash)
	}
	var memoFollows []*MemoFollow
	result := query.Scan(&memoFollows)
	if result.Error != nil {
		return nil, jerr.Get("error running follower query", result.Error)
	}
	return memoFollows, nil
}

func GetFollowingForPkHash(followPkHash []byte, offset int) ([]*MemoFollow, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	sql := "" +
		"SELECT " +
		"	memo_follows.* " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE follow_pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id) " +
		"WHERE unfollow = 0 "
	var query *gorm.DB
	if offset >= 0 {
		sql += "LIMIT ?,25"
		query = db.Raw(sql, followPkHash, offset)
	} else {
		query = db.Raw(sql, followPkHash)
	}
	var memoFollows []*MemoFollow
	result := query.Scan(&memoFollows)
	if result.Error != nil {
		return nil, jerr.Get("error running following query", result.Error)
	}
	return memoFollows, nil
}

func GetFollowingCountForPkHash(pkHash []byte) (uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	sql := "" +
		"SELECT " +
		"	COALESCE(SUM(IF(unfollow=0, 1, 0)), 0) AS following " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id)"
	query := db.Raw(sql, pkHash)
	var cnt uint
	row := query.Row()
	err = row.Scan(&cnt)
	if err != nil {
		return 0, jerr.Get("error running following count query", err)
	}
	return cnt, nil
}

func GetFollowerCountForPkHash(followPkHash []byte) (uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	sql := "" +
		"SELECT " +
		"	COALESCE(SUM(IF(unfollow=0, 1, 0)), 0) AS followers " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE follow_pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id)"
	query := db.Raw(sql, followPkHash)
	var cnt uint
	row := query.Row()
	err = row.Scan(&cnt)
	if err != nil {
		return 0, jerr.Get("error running follower count query", err)
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

func IsFollowing(followerPkHash []byte, followingPkHash []byte) (bool, error) {
	if len(followerPkHash) == 0 {
		return false, nil
	}
	db, err := getDb()
	if err != nil {
		return false, jerr.Get("error getting db", err)
	}
	sql := "" +
		"SELECT " +
		"	COALESCE(unfollow, 1) AS is_following " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE pk_hash = ? AND follow_pk_hash = ?" +
		") sq ON (sq.id = memo_follows.id)"
	query := db.Raw(sql, followerPkHash, followingPkHash)
	var cnt uint
	row := query.Row()
	err = row.Scan(&cnt)
	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return false, nil
		}
		return false, jerr.Get("error is follower query", err)
	}
	return cnt == 0, nil
}
