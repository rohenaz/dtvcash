package topics

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var indexRoute = web.Route{
	Pattern: res.UrlTopics,
	Handler: func(r *web.Response) {
		preHandler(r)
		topics, err := db.GetUniqueTopics()
		if err != nil {
			r.Error(jerr.Get("error getting topics from db", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Topics"] = topics
		r.Render()
	},
}
