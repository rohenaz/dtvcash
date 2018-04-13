package server

import (
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/web"
)

var protocolRoute = web.Route{
	Pattern: res.UrlProtocol,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Protocol"
		r.Render()
	},
}
