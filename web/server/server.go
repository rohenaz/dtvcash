package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/res"
	"git.jasonc.me/main/memo/web/server/key"
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
	r.Helper["BaseUrl"] = getBaseUrl(r)
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

func getBaseUrl(r *web.Response) string {
	baseUrl := r.Request.GetHeader("AppPath")
	if baseUrl == "" {
		baseUrl = "/"
	}
	return baseUrl
}

func getUrlWithBaseUrl(url string, r *web.Response) string {
	baseUrl := getBaseUrl(r)
	baseUrl = baseUrl[:len(baseUrl)-1]
	return baseUrl + url
}

func Run(sessionCookieInsecure bool) {
	// Start bitcoin node
	res.BitcoinNode.Address = res.BitcoinPeerAddress
	res.BitcoinNode.Start()

	// Start web server
	ws := web.Server{
		CookiePrefix:   "memo",
		InsecureCookie: sessionCookieInsecure,
		IsLoggedIn:     isLoggedIn,
		Port:           8261,
		PreHandler:     preHandler,
		Routes: append([]web.Route{
			indexRoute,
			loginRoute,
			loginSubmitRoute,
			logoutRoute,
			signupRoute,
			signupSubmitRoute,
		}, key.GetRoutes()...),
		StaticFilesDir: "web/public",
		TemplatesDir:   "web/templates",
		UseSessions:    true,
	}
	err := ws.Run()
	if err != nil {
		log.Fatal(err)
	}
}
