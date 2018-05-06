package topics

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
	"github.com/memocash/memo/app/watcher"
	"net/http"
	"net/url"
)

var socketRoute = web.Route{
	Pattern: res.UrlTopicsSocket,
	Handler: func(r *web.Response) {
		topicRaw := r.Request.GetUrlParameter("topic")
		unescaped, err := url.QueryUnescape(topicRaw)
		if err != nil {
			r.Error(jerr.Get("error unescaping topic", err), http.StatusUnprocessableEntity)
			return
		}
		lastPostId := r.Request.GetUrlParameterUInt("lastPostId")
		lastLikeId := r.Request.GetUrlParameterUInt("lastLikeId")
		socket, err := r.GetWebSocket()
		if err != nil {
			r.Error(jerr.Get("error getting socket", err), http.StatusUnprocessableEntity)
			return
		}
		err = watcher.RegisterSocket(socket, unescaped, lastPostId, lastLikeId)
		if err != nil {
			r.Error(jerr.Get("error writing to socket", err), http.StatusInternalServerError)
			return
		}
	},
}
