package profile

import (
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

type Post struct {
	Name string
	Memo *db.MemoPost
}

func GetPostsForHashes(pkHashes [][]byte) ([]*Post, error) {
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
		name, ok := names[string(dbPost.PkHash)]
		if ok {
			post.Name = name
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsForHash(pkHash []byte) ([]*Post, error) {
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
		posts = append(posts, &Post{
			Name: setName.Name,
			Memo: dbPost,
		})
	}
	return posts, nil
}
