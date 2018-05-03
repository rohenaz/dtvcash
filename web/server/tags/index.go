package tags

import (
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/web"
)

var indexRoute = web.Route{
	Pattern: res.UrlTags,
	Handler: func(r *web.Response) {
		preHandler(r)
		r.Render()
	},
}
