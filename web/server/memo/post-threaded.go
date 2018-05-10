package memo

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var postThreadedRoute = web.Route{
	Pattern: res.UrlMemoPostThreaded + "/" + urlTxHash.UrlPart(),
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		post, err := getPostWithThreads(r, txHashString, offset)
		if err != nil {
			r.Error(jerr.Get("error getting post with threads", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Post"] = post
		r.Helper["Offset"] = 0
		r.Helper["Title"] = fmt.Sprintf("Memo - Post by %s", post.Name)
		if post.Name == "" {
			r.Helper["Title"] = fmt.Sprintf("Memo - Post by %.6s", post.Memo.GetAddressString())
		}
		r.Helper["Description"] = post.Memo.Message
		r.RenderTemplate(res.TmplMemoPostThreaded)
	},
}

var postThreadedAjaxRoute = web.Route{
	Pattern: res.UrlMemoPostThreadedAjax,
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		txHashString := r.Request.GetUrlParameter("txHash")
		post, err := getPostWithThreads(r, txHashString, offset)
		if err != nil {
			r.Error(jerr.Get("error getting post with threads", err), http.StatusInternalServerError)
			return
		}
		for _, reply := range post.Replies {
			r.Helper["Post"] = reply
			r.RenderTemplate(res.TmplSnippetsPostThreaded)
		}
		if len(post.Replies) == 25 {
			r.Helper["Post"] = post
			r.Helper["Offset"] = offset
			r.RenderTemplate(res.TmplSnippetsPostThreadedLoadMore)
		}
	},
}

func getPostWithThreads(r *web.Response, txHashString string, offset int) (*profile.Post, error) {
	txHash, err := chainhash.NewHashFromStr(txHashString)
	if err != nil {
		return nil, jerr.Get("error getting transaction hash", err)
	}
	var pkHash []byte
	if auth.IsLoggedIn(r.Session.CookieId) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			return nil, jerr.Get("error getting session user", err)
		}
		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			return nil, jerr.Get("error getting key for user", err)
		}
		pkHash = key.PkHash
	}
	post, err := profile.GetPostByTxHash(txHash.CloneBytes(), pkHash, uint(offset))
	if err != nil {
		return nil, jerr.Get("error getting post", err)
	}
	allPosts := append(post.Replies, post)
	needsReplies := post.Replies
	for len(needsReplies) != 0 {
		needsRepliesPost := needsReplies[0]
		needsReplies = needsReplies[1:]
		err = profile.AttachRepliesToPost(needsRepliesPost, 0)
		if err != nil {
			return nil, jerr.Get("error attaching replies to reply", err)
		}
		allPosts = append(allPosts, needsRepliesPost.Replies...)
		needsReplies = append(needsReplies, needsRepliesPost.Replies...)
		if len(allPosts) > 250 {
			jerr.New("Nested replies over 250, breaking out!").Print()
			break
		}
	}
	err = profile.AttachLikesToPosts(allPosts)
	if err != nil {
		return nil, jerr.Get("error attaching likes to posts", err)
	}
	if len(pkHash) > 0 {
		err = profile.AttachReputationToPosts(allPosts)
		if err != nil {
			return nil, jerr.Get("error attaching reputation to posts", err)
		}
	}
	return post, nil
}
