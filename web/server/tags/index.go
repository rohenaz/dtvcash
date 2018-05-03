package tags

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var indexRoute = web.Route{
	Pattern: res.UrlTags,
	Handler: func(r *web.Response) {
		preHandler(r)
		tags, err := db.GetUniqueTags()
		if err != nil {
			r.Error(jerr.Get("error getting tags from db", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Tags"] = tags
		r.Render()
	},
}
