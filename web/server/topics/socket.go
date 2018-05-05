package topics

import (
	"git.jasonc.me/main/memo/app/html-parser"
	"git.jasonc.me/main/memo/app/res"
	"git.jasonc.me/main/memo/app/watcher"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var socketRoute = web.Route{
	Pattern: res.UrlTopicsSocket,
	Handler: func(r *web.Response) {
		topicRaw := r.Request.GetUrlParameter("topic")
		lastPostId := r.Request.GetUrlParameterUInt("lastPostId")
		lastLikeId := r.Request.GetUrlParameterUInt("lastLikeId")
		topic := html_parser.EscapeWithEmojis(topicRaw)
		socket, err := r.GetWebSocket()
		if err != nil {
			r.Error(jerr.Get("error getting socket", err), http.StatusUnprocessableEntity)
			return
		}
		err = watcher.RegisterSocket(socket, topic, lastPostId, lastLikeId)
		if err != nil {
			r.Error(jerr.Get("error writing to socket", err), http.StatusInternalServerError)
			return
		}
	},
}
