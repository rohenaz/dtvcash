package profile

import (
	"bytes"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"regexp"
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

func (p Post) GetMessage() string {
	msg := p.Memo.Message
	var re = regexp.MustCompile(`(http[s]?://[^ ]*)`)
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
	var name = ""
	if setName != nil {
		name = setName.Name
	}
	post := &Post{
		Name:       name,
		Memo:       memoPost,
		SelfPkHash: selfPkHash,
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
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
		}
		posts = append(posts, post)
	}
	return posts, nil
}
