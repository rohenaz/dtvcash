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

var postRoute = web.Route{
	Pattern:    res.UrlMemoPost + "/" + urlTxHash.UrlPart(),
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		post, err := getPostWithThreads(r, txHashString, offset)
		if err != nil {
			if db.IsRecordNotFoundError(err) {
				r.Error(jerr.Get("error post not found", err), http.StatusNotFound)
				r.RenderTemplate(res.UrlNotFound)
				return
			}
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

var postAjaxRoute = web.Route{
	Pattern:    res.UrlMemoPostAjax + "/" + urlTxHash.UrlPart(),
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			 r.Error(jerr.Get("error getting transaction hash", err), http.StatusUnprocessableEntity)
			 return
		}
		var pkHash []byte
		if auth.IsLoggedIn(r.Session.CookieId) {
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
				return
			}
			key, err := db.GetKeyForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
				return
			}
			pkHash = key.PkHash
		}
		post, err := profile.GetPostByTxHash(txHash.CloneBytes(), pkHash)
		if err != nil {
			r.Error(jerr.Get("error getting post", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachParentToPosts([]*profile.Post{post})
		if err != nil {
			r.Error(jerr.Get("error attaching parent to post", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts([]*profile.Post{post})
		if err != nil {
			r.Error(jerr.Get("error attaching likes to post", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Post"] = post
		r.Helper["Offset"] = 0
		r.RenderTemplate(res.TmplSnippetsPost)
	},
}

var postThreadedAjaxRoute = web.Route{
	Pattern:    res.UrlMemoPostThreadedAjax + "/" + urlTxHash.UrlPart(),
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
		r.RenderTemplate(res.TmplSnippetsPostThreaded)
	},
}
