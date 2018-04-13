package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/res"
	auth2 "git.jasonc.me/main/memo/web/server/auth"
	"git.jasonc.me/main/memo/web/server/key"
	"git.jasonc.me/main/memo/web/server/memo"
	"git.jasonc.me/main/memo/web/server/profile"
	"github.com/jchavannes/jgo/web"
	"log"
	"net/http"
)

var UseMinJS bool

func isLoggedIn(r *web.Response) bool {
	if ! auth.IsLoggedIn(r.Session.CookieId) {
		r.SetResponseCode(http.StatusUnauthorized)
		return false
	}
	return true
}

func preHandler(r *web.Response) {
	r.Helper["Title"] = "Memo"
	r.Helper["BaseUrl"] = res.GetBaseUrl(r)
	if auth.IsLoggedIn(r.Session.CookieId) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.Helper["Username"] = user.Username
	}
	if UseMinJS {
		r.Helper["jsFiles"] = res.GetMinJsFiles()
	} else {
		r.Helper["jsFiles"] = res.GetResJsFiles()
	}
	r.Helper["cssFiles"] = res.CssFiles
}

func Run(sessionCookieInsecure bool) {
	go func() {
		// Start bitcoin node
		node.BitcoinNode.NetAddress = node.BitcoinPeerAddress
		node.BitcoinNode.SetKeys()
		node.BitcoinNode.Start()
	}()

	// Start web server
	ws := web.Server{
		CookiePrefix:   "memo",
		InsecureCookie: sessionCookieInsecure,
		IsLoggedIn:     isLoggedIn,
		Port:           8261,
		PreHandler:     preHandler,
		Routes: web.Routes(
			[]web.Route{
				indexRoute,
				protocolRoute,
				testsRoute,
			},
			key.GetRoutes(),
			auth2.GetRoutes(),
			memo.GetRoutes(),
			profile.GetRoutes(),
		),
		StaticFilesDir: "web/public",
		TemplatesDir:   "web/templates",
		UseSessions:    true,
	}
	err := ws.Run()
	if err != nil {
		log.Fatal(err)
	}
}
