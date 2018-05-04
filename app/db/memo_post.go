package db

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/script"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"html"
	"sort"
	"time"
)

const (
	PreloadMemoPostParent = "Parent"
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
	Topic        string      `gorm:"index:tag"`
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

func (m MemoPost) GetParentTransactionHashString() string {
	hash, err := chainhash.NewHash(m.ParentTxHash)
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

func GetMemoPostById(id uint) (*MemoPost, error) {
	var memoPost MemoPost
	err := find(&memoPost, MemoPost{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoPost, nil
}

func GetPostReplyCount(txHash []byte) (uint, error) {
	cnt, err := count(MemoPost{
		ParentTxHash: txHash,
	})
	if err != nil {
		return 0, jerr.Get("error running count query", err)
	}
	return cnt, nil
}

func GetPostReplies(txHash []byte, offset uint) ([]*MemoPost, error) {
	var posts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}

	query := db.
		Table("memo_posts").
		Preload(BlockTable).
		Select("memo_posts.*, COUNT(DISTINCT memo_likes.pk_hash) AS count").
		Joins("LEFT OUTER JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Joins("LEFT OUTER JOIN memo_likes ON (memo_posts.tx_hash = memo_likes.like_tx_hash)").
		Group("memo_posts.id").
		Order("count DESC, memo_posts.id DESC").
		Limit(25).
		Offset(offset)

	result := query.Find(&posts, MemoPost{
		ParentTxHash: txHash,
	})
	if result.Error != nil {
		return nil, jerr.Get("error finding post replies", result.Error)
	}
	return posts, nil
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
		return txns[i].Id > txns[j].Id
	}
	if txns[i].Block == nil {
		return true
	}
	if txns[j].Block == nil {
		return false
	}
	if txns[i].Block.Height == txns[j].Block.Height {
		return txns[i].Id > txns[j].Id
	}
	return txns[i].Block.Height > txns[j].Block.Height
}

type memoPostSortByDateAsc []*MemoPost

func (txns memoPostSortByDateAsc) Len() int      { return len(txns) }
func (txns memoPostSortByDateAsc) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoPostSortByDateAsc) Less(i, j int) bool {
	if bytes.Equal(txns[i].ParentHash, txns[j].TxHash) {
		return false
	}
	if bytes.Equal(txns[i].TxHash, txns[j].ParentHash) {
		return true
	}
	if txns[i].Block == nil && txns[j].Block == nil {
		return txns[i].Id < txns[j].Id
	}
	if txns[i].Block == nil {
		return false
	}
	if txns[j].Block == nil {
		return true
	}
	if txns[i].Block.Height == txns[j].Block.Height {
		return txns[i].Id < txns[j].Id
	}
	return txns[i].Block.Height < txns[j].Block.Height
}

func GetPostsForPkHashes(pkHashes [][]byte, offset uint) ([]*MemoPost, error) {
	if len(pkHashes) == 0 {
		return nil, nil
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.
		Limit(25).
		Offset(offset).
		Preload(BlockTable).
		Where("pk_hash in (?)", pkHashes).
		Order("id DESC").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetPostsFeedForPkHash(pkHash []byte, offset uint) ([]*MemoPost, error) {
	var memoPosts []*MemoPost
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
	result := db.
		Limit(25).
		Offset(offset).
		Preload(BlockTable).
		Joins("JOIN ("+joinSelect+") fsq ON (memo_posts.pk_hash = fsq.follow_pk_hash)", pkHash).
		Order("id DESC").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetPostsForPkHash(pkHash []byte, offset uint) ([]*MemoPost, error) {
	if len(pkHash) == 0 {
		return nil, nil
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Order("id DESC").
		Limit(25).
		Offset(offset)
	result := query.Find(&memoPosts, &MemoPost{
		PkHash: pkHash,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
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
	defer rows.Close()
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
	result := db.
		Limit(25).
		Offset(offset).
		Order("id DESC").
		Where("address != ?", "15wck91ux5557CK4Wb42Vj6pqwdosgwd7H").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	sort.Sort(memoPostSortByDate(memoPosts))
	return memoPosts, nil
}

func GetRecentPostsForTopic(topic string, lastPostId uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoPosts []*MemoPost
	result := db.
		Where("id > ?", lastPostId).
		Order("id ASC").
		Find(&memoPosts, MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error running recent topic post query", result.Error)
	}
	return memoPosts, nil
}

func GetTopPosts(offset uint, timeStart time.Time, timeEnd time.Time) ([]*MemoPost, error) {
	topLikeTxHashes, err := GetRecentTopLikedTxHashes(offset, timeStart, timeEnd)
	if err != nil {
		return nil, jerr.Get("error getting top liked tx hashes", err)
	}
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.Where("tx_hash IN (?)", topLikeTxHashes).Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	var sortedPosts []*MemoPost
	for _, txHash := range topLikeTxHashes {
		for _, memoPost := range memoPosts {
			if bytes.Equal(memoPost.TxHash, txHash) {
				sortedPosts = append(sortedPosts, memoPost)
			}
		}
	}
	return sortedPosts, nil
}

func GetPersonalizedTopPosts(selfPkHash []byte, offset uint, timeStart time.Time, timeEnd time.Time) ([]*MemoPost, error) {
	topLikeTxHashes, err := GetPersonalizedRecentTopLikedTxHashes(selfPkHash, offset, timeStart, timeEnd)
	if err != nil {
		return nil, jerr.Get("error getting top liked tx hashes", err)
	}
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.Where("tx_hash IN (?)", topLikeTxHashes).Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	var sortedPosts []*MemoPost
	for _, txHash := range topLikeTxHashes {
		for _, memoPost := range memoPosts {
			if bytes.Equal(memoPost.TxHash, txHash) {
				sortedPosts = append(sortedPosts, memoPost)
			}
		}
	}
	return sortedPosts, nil
}

func GetCountMemoPosts() (uint, error) {
	cnt, err := count(&MemoPost{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}

type Topic struct {
	Name       string
	RecentTime time.Time
	Count      int
}

func GetUniqueTopics(offset uint) ([]*Topic, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	rows, err := db.
		Table("memo_posts").
		Select("topic, MAX(memo_posts.created_at) AS max_time, COUNT(*)").
		Joins("LEFT OUTER JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Where("topic IS NOT NULL AND topic != ''").
		Group("topic").
		Order("max_time DESC").
		Limit(25).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, jerr.Get("error getting distinct topics", err)
	}
	defer rows.Close()
	var topics []*Topic
	for rows.Next() {
		var topic Topic

		err := rows.Scan(&topic.Name, &topic.RecentTime, &topic.Count)
		if err != nil {
			return nil, jerr.Get("error scanning row with topic", err)
		}
		topics = append(topics, &topic)
	}
	return topics, nil
}

func GetPostsForTopic(topic string, offset uint) ([]*MemoPost, error) {
	if len(topic) == 0 {
		return nil, jerr.New("empty topic")
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Order("id DESC").
		Limit(26).
		Offset(offset)
	result := query.Find(&memoPosts, &MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	sort.Sort(memoPostSortByDateAsc(memoPosts))
	return memoPosts, nil
}

func GetOlderPostsForTopic(topic string, firstPostId uint) ([]*MemoPost, error) {
	if len(topic) == 0 {
		return nil, jerr.New("empty topic")
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Where("id < ?", firstPostId).
		Order("id DESC").
		Limit(26)
	result := query.Find(&memoPosts, &MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	sort.Sort(memoPostSortByDateAsc(memoPosts))
	return memoPosts, nil
}
