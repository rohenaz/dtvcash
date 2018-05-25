package topics

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/html-parser"
	"github.com/rohenaz/dtvcash/app/profile"
	"github.com/rohenaz/dtvcash/app/res"
	"net/http"
)

var postsMoreRoute = web.Route{
	Pattern: res.UrlTopicsMorePosts,
	Handler: func(r *web.Response) {
		firstPostId := r.Request.GetUrlParameterUInt("firstPostId")
		topicRaw := r.Request.GetUrlParameter("topic")
		topic := html_parser.EscapeWithEmojis(topicRaw)
		var userPkHash []byte
		var userId uint
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
			userPkHash = key.PkHash
			userId = user.Id
		}
		posts, err := profile.GetOlderPostsForTopic(topic, userPkHash, firstPostId)
		if err != nil {
			r.Error(jerr.Get("error getting posts", err), http.StatusInternalServerError)
			return
		}
		if len(posts) == 0 {
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(posts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachReplyCountToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching reply counts to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.SetShowMediaForPosts(posts, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Posts"] = posts
		r.Helper["FirstPostId"] = posts[0].Memo.Id
		r.Render()
	},
}

var postAjaxRoute = web.Route{
	Pattern: res.UrlTopicsPostAjax + "/" + urlTxHash.UrlPart(),
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		var pkHash []byte
		var userId uint
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
			userId = user.Id
		}
		post, err := profile.GetPostByTxHashWithReplies(txHash.CloneBytes(), pkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting post", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachParentToPosts([]*profile.Post{post})
		if err != nil {
			r.Error(jerr.Get("error attaching parent to post", err), http.StatusInternalServerError)
			return
		}
		if len(pkHash) > 0 {
			err = profile.AttachReputationToPosts([]*profile.Post{post})
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to post", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.AttachLikesToPosts([]*profile.Post{post})
		if err != nil {
			r.Error(jerr.Get("error attaching likes to post", err), http.StatusInternalServerError)
			return
		}
		err = profile.SetShowMediaForPosts([]*profile.Post{post}, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Post"] = post
		r.RenderTemplate(res.TmplTopicPost)
	},
}
