package profile

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"time"
)

type Like struct {
	Id         uint
	Amount     int64
	Name       string
	PkHash     []byte
	Timestamp  time.Time
	TxnHash    []byte
	PostTxHash []byte
}

func (l *Like) GetAddressString() string {
	return wallet.GetAddressFromPkHash(l.PkHash).GetEncoded()
}

func (l *Like) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(l.TxnHash)
	if err != nil {
		return ""
	}
	return hash.String()
}

func (l *Like) GetPostTransactionHashString() string {
	hash, err := chainhash.NewHash(l.PostTxHash)
	if err != nil {
		return ""
	}
	return hash.String()
}

func (l *Like) GetTimeString(timezone string) string {
	if ! l.Timestamp.IsZero() {
		timeLayout := "2006-01-02 15:04:05 MST"
		if len(timezone) > 0 {
			displayLocation, err := time.LoadLocation(timezone)
			if err != nil {
				jerr.Get("error loading location", err).Print()
				return l.Timestamp.Format(timeLayout)
			}
			return l.Timestamp.In(displayLocation).Format(timeLayout)
		} else {
			return l.Timestamp.Format(timeLayout)
		}
	}
	return "Unconfirmed"
}

func AttachLikesToPosts(posts []*Post) error {
	for _, post := range posts {
		memoLikes, err := db.GetMemoLikesForTxnHash(post.Memo.TxHash)
		if err != nil {
			return jerr.Get("error getting likes for post", err)
		}
		var likes []*Like
		for _, memoLike := range memoLikes {
			var like = &Like{
				Id:      memoLike.Id,
				PkHash:  memoLike.PkHash,
				TxnHash: memoLike.TxHash,
			}
			if memoLike.Block != nil {
				like.Timestamp = memoLike.Block.Timestamp
			}
			if bytes.Equal(memoLike.TipPkHash, post.Memo.PkHash) {
				like.Amount = memoLike.TipAmount
			}
			setName, err := db.GetNameForPkHash(memoLike.PkHash)
			if err != nil && ! db.IsRecordNotFoundError(err) {
				return jerr.Get("error getting memo name", err)
			}
			if setName != nil {
				like.Name = setName.Name
			}
			if bytes.Equal(like.PkHash, post.SelfPkHash) {
				post.HasLiked = true
			}
			likes = append(likes, like)
		}
		post.Likes = likes
	}
	return nil
}

func AttachReputationToPosts(posts []*Post) error {
	for _, post := range posts {
		reputation, err := GetReputation(post.SelfPkHash, post.Memo.PkHash)
		if err != nil {
			return jerr.Get("error getting reputation", err)
		}
		post.Reputation = reputation
	}
	return nil
}

func GetLikesForPkHash(pkHash []byte) ([]*Like, error) {
	memoLikes, err := db.GetMemoLikesForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting memo likes from db", err)
	}
	var likes []*Like
	for _, memoLike := range memoLikes {
		var like = Like{
			PkHash:     pkHash,
			TxnHash:    memoLike.TxHash,
			PostTxHash: memoLike.LikeTxHash,
		}
		memoPost, err := db.GetMemoPost(memoLike.LikeTxHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting transaction from db", err)
		}
		if memoPost != nil && bytes.Equal(memoLike.TipPkHash, memoPost.PkHash) {
			like.Amount = memoLike.TipAmount
		}
		if memoLike.Block != nil {
			like.Timestamp = memoLike.Block.Timestamp
		}
		if memoPost != nil {
			setName, err := db.GetNameForPkHash(memoPost.PkHash)
			if err != nil && ! db.IsRecordNotFoundError(err) {
				return nil, jerr.Get("error getting memo name", err)
			}
			if setName != nil {
				like.Name = setName.Name
			}
		}
		likes = append(likes, &like)
	}
	return likes, nil
}
