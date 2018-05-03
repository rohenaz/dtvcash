package tags

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"time"
)

var socketRoute = web.Route{
	Pattern: res.UrlTagsSocket,
	Handler: func(r *web.Response) {
		tagRaw := r.Request.GetUrlParameter("tag")
		socket, err := r.GetWebSocket()
		if err != nil {
			r.Error(jerr.Get("error getting socket", err), http.StatusUnprocessableEntity)
			return
		}
		var lastPostId uint
		for i := 0; i < 1e6; i++ {
			recentPost, err := db.GetRecentPostForTag(tagRaw)
			if err != nil {
				r.Error(jerr.Get("error getting recent post for tag", err), http.StatusInternalServerError)
				return
			}
			if lastPostId == 0 {
				lastPostId = recentPost.Id
			}
			if recentPost.Id != lastPostId {
				fmt.Println("Found new post!")
				lastPostId = recentPost.Id
				socket.WriteJSON(recentPost.GetTransactionHashString())
			}
			time.Sleep(250 * time.Millisecond)
		}
	},
}
