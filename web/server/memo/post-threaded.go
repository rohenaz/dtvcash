package memo

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/profile"
	"github.com/rohenaz/dtvcash/app/res"
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

var postMoreThreadedAjaxRoute = web.Route{
	Pattern: res.UrlMemoPostMoreThreadedAjax,
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
	var userId uint
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
		userId = user.Id
	}
	post, err := profile.GetPostByTxHashWithReplies(txHash.CloneBytes(), pkHash, uint(offset))
	if err != nil {
		return nil, jerr.Get("error getting post", err)
	}
	err = profile.AttachParentToPosts([]*profile.Post{post})
	if err != nil {
		return nil, jerr.Get("error attaching parent to post", err)
	}
	allPosts := []*profile.Post{post}
	needsReplies := post.Replies
	for len(needsReplies) != 0 {
		needsRepliesPost := needsReplies[0]
		needsReplies = needsReplies[1:]
		err = profile.AttachRepliesToPost(needsRepliesPost, 0)
		if err != nil {
			return nil, jerr.Get("error attaching replies to reply", err)
		}
		allPosts = append(allPosts, needsRepliesPost)
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
	err = profile.AttachPollsToPosts(allPosts)
	if err != nil {
		return nil, jerr.Get("error attaching polls to posts", err)
	}
	err = profile.SetShowMediaForPosts(allPosts, userId)
	if err != nil {
		return nil, jerr.Get("error setting show media for posts", err)
	}
	return post, nil
}
