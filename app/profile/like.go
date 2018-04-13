package profile

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type Like struct {
	Amount    int64
	Name      string
	PkHash    []byte
	Timestamp time.Time
	TxnHash   []byte
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

func (l *Like) GetTimeString() string {
	if ! l.Timestamp.IsZero() {
		return l.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func AttachLikesToPosts(posts []*Post) error {
	for _, post := range posts {
		memoLikes, err := db.GetLikesForTxnHash(post.Memo.TxHash)
		if err != nil {
			return jerr.Get("error getting likes for post", err)
		}
		var likes []*Like
		for _, memoLike := range memoLikes {
			like := &Like{
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
			likes = append(likes, like)
		}
		post.Likes = likes
	}
	return nil
}
