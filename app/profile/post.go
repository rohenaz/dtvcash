package profile

import (
	"bytes"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"regexp"
	"time"
)

type Post struct {
	Name       string
	Memo       *db.MemoPost
	Parent     *Post
	Likes      []*Like
	SelfPkHash []byte
	ReplyCount uint
	Replies    []*Post
}

func (p Post) IsSelf() bool {
	if len(p.SelfPkHash) == 0 {
		return false
	}
	return bytes.Equal(p.SelfPkHash, p.Memo.PkHash)
}

func (p Post) IsLikable() bool {
	return len(p.SelfPkHash) > 0
}

func (p Post) GetTotalTip() int64 {
	var totalTip int64
	for _, like := range p.Likes {
		totalTip += like.Amount
	}
	return totalTip
}

func (p Post) GetMessage() string {
	msg := p.Memo.Message
	var re = regexp.MustCompile(`(http[s]?://[^\s]*)`)
	s := re.ReplaceAllString(msg, `<a href="$1" target="_blank">$1</a>`)
	return s
}

func GetPostsForHashes(pkHashes [][]byte, selfPkHash []byte, offset uint) ([]*Post, error) {
	dbPosts, err := db.GetPostsForPkHashes(pkHashes, offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var foundPkHashes [][]byte
	for _, dbPost := range dbPosts {
		for _, foundPkHash := range foundPkHashes {
			if bytes.Equal(foundPkHash, dbPost.PkHash) {
				continue
			}
		}
		foundPkHashes = append(foundPkHashes, dbPost.PkHash)
	}
	names := make(map[string]string)
	for _, pkHash := range foundPkHashes {
		setName, err := db.GetNameForPkHash(pkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName == nil {
			continue
		}
		names[string(pkHash)] = setName.Name
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		name, ok := names[string(dbPost.PkHash)]
		if ok {
			post.Name = name
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsForHash(pkHash []byte, selfPkHash []byte) ([]*Post, error) {
	var name = ""
	setName, err := db.GetNameForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	if setName != nil {
		name = setName.Name
	}
	dbPosts, err := db.GetPostsForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByTxHash(txHash []byte, selfPkHash []byte) (*Post, error) {
	memoPost, err := db.GetMemoPost(txHash)
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	var parent *Post
	if len(memoPost.ParentTxHash) > 0 {
		parentPost, err := db.GetMemoPost(memoPost.ParentTxHash)
		if err != nil {
			return nil, jerr.Get("error getting memo post parent", err)
		}
		setName, err := db.GetNameForPkHash(parentPost.PkHash)
		if err != nil {
			return nil, jerr.Get("error getting name for reply hash", err)
		}
		var name = ""
		if setName != nil {
			name = setName.Name
		}
		parent = &Post{
			Name:       name,
			Memo:       parentPost,
			SelfPkHash: selfPkHash,
		}
	}
	replies, err := db.GetPostReplies(txHash)
	if err != nil {
		return nil, jerr.Get("error getting post replies", err)
	}
	var replyPosts []*Post
	for _, reply := range replies {
		setName, err := db.GetNameForPkHash(reply.PkHash)
		if err != nil {
			return nil, jerr.Get("error getting name for reply hash", err)
		}
		var name = ""
		if setName != nil {
			name = setName.Name
		}
		cnt, err := db.GetPostReplyCount(reply.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		replyPosts = append(replyPosts, &Post{
			Name:       name,
			Memo:       reply,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		})
	}
	setName, err := db.GetNameForPkHash(memoPost.PkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	var name = ""
	if setName != nil {
		name = setName.Name
	}
	post := &Post{
		Name:       name,
		Memo:       memoPost,
		Parent:     parent,
		SelfPkHash: selfPkHash,
		Replies:    replyPosts,
		ReplyCount: uint(len(replies)),
	}
	return post, nil
}

func GetRecentPosts(selfPkHash []byte, offset uint) ([]*Post, error) {
	dbPosts, err := db.GetRecentPosts(offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		var name string
		setName, err := db.GetNameForPkHash(dbPost.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName != nil {
			name = setName.Name
		}
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p Post) GetTimeString(timezone string) string {
	if p.Memo.BlockId != 0 {
		if p.Memo.Block != nil {
			timeLayout := "2006-01-02 15:04:05 MST"
			if len(timezone) > 0 {
				displayLocation, err := time.LoadLocation(timezone)
				if err != nil {
					jerr.Get("error finding location", err).Print()
					return p.Memo.Block.Timestamp.Format(timeLayout)
				}
				return p.Memo.Block.Timestamp.In(displayLocation).Format(timeLayout)
			} else {
				return p.Memo.Block.Timestamp.Format(timeLayout)
			}
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}