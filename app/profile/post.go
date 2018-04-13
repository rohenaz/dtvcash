package profile

import (
	"bytes"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

type Post struct {
	Name  string
	Memo  *db.MemoPost
	Likes []*db.MemoLike
	Self  bool
}

func (p Post) IsSelf() bool {
	return p.Self
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
			Memo: dbPost,
		}
		if bytes.Equal(dbPost.PkHash, selfPkHash) {
			post.Self = true
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
	setName, err := db.GetNameForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	dbPosts, err := db.GetPostsForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		post := &Post{
			Name: setName.Name,
			Memo: dbPost,
		}
		if bytes.Equal(dbPost.PkHash, selfPkHash) {
			post.Self = true
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
		Name: setName.Name,
		Memo: memoPost,
	}
	if bytes.Equal(post.Memo.PkHash, selfPkHash) {
		post.Self = true
	}
	return post, nil
}

func AttachLikesToPosts(posts []*Post) error {
	for _, post := range posts {
		memoLikes, err := db.GetLikesForTxnHash(post.Memo.TxHash)
		if err != nil {
			return jerr.Get("error getting likes for post", err)
		}
		post.Likes = memoLikes
	}
	return nil
}
