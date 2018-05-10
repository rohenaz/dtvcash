package memo

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
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
