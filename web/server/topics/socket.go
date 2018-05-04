package topics

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/html-parser"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"time"
)

var socketRoute = web.Route{
	Pattern: res.UrlTopicsSocket,
	Handler: func(r *web.Response) {
		topicRaw := r.Request.GetUrlParameter("topic")
		lastPostId := r.Request.GetUrlParameterUInt("lastPostId")
		topic := html_parser.EscapeWithEmojis(topicRaw)
		socket, err := r.GetWebSocket()
		if err != nil {
			r.Error(jerr.Get("error getting socket", err), http.StatusUnprocessableEntity)
			return
		}
		for i := 0; i < 1e6; i++ {
			recentPosts, err := db.GetRecentPostsForTopic(topic, lastPostId)
			if err != nil && !db.IsRecordNotFoundError(err) {
				r.Error(jerr.Get("error getting recent post for topic", err), http.StatusInternalServerError)
				return
			}
			if len(recentPosts) > 0 {
				fmt.Println("Found new post(s)!")
				for _, recentPost := range recentPosts {
					lastPostId = recentPost.Id
					socket.WriteJSON(recentPost.GetTransactionHashString())
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	},
}
