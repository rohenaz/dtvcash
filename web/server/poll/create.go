package poll

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var createRoute = web.Route{
	Pattern:    res.UrlPollCreate,
	Handler: func(r *web.Response) {
		r.Render()
	},
}
