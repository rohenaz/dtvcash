package profile

import (
	"bytes"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

type Post struct {
	Name       string
	Memo       *db.MemoPost
	Likes      []*Like
	SelfPkHash []byte
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

func GetPostsForHashes(pkHashes [][]byte, selfPkHash []byte) ([]*Post, error) {
	names := make(map[string]string)
	for _, pkHash := range pkHashes {
		setName, err := db.GetNameForPkHash(pkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if db.IsRecordNotFoundError(err) {
			continue
		}
		names[string(pkHash)] = setName.Name
	}
	dbPosts, err := db.GetPostsForPkHashes(pkHashes)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		post := &Post{
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
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
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByTxHash(txHash []byte, selfPkHash []byte) (*Post, error) {
	memoPost, err := db.GetMemoPost(txHash)
	if err != nil {
		return nil, jerr.Get("error getting post", err)
	}
	setName, err := db.GetNameForPkHash(memoPost.PkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	post := &Post{
		Name:       setName.Name,
		Memo:       memoPost,
		SelfPkHash: selfPkHash,
	}
	return post, nil
}
