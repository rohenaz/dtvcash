package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/web"
	"log"
	"net/http"
)

const (
	UrlIndex        = "/"
	UrlSignup       = "/signup"
	UrlSignupSubmit = "/signup-submit"
	UrlLogin        = "/login"
	UrlLoginSubmit  = "/login-submit"
	UrlLogout       = "/logout"
)
const (
	UrlKeyView                = "/key"
	UrlKeyLoad                = "/key/load"
	UrlKeyImport              = "/key/import"
	UrlKeyImportSubmit        = "/key/import-submit"
	UrlKeyCreate              = "/key/create"
	UrlCreatePrivateKeySubmit = "/key/create-submit"
	UrlKeyDeleteSubmit        = "/key/delete-submit"
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

var urlId = web.UrlParam{
	Id:   "id",
	Type: web.UrlParamInteger,
}

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var bitcoinNode node.Node

func Run(sessionCookieInsecure bool) {
	// Start bitcoin node
	bitcoinNode.Address = BitcoinPeerAddress
	bitcoinNode.Start()

	// Start web server
	ws := web.Server{
		CookiePrefix:   "memo",
		InsecureCookie: sessionCookieInsecure,
		IsLoggedIn:     isLoggedIn,
		Port:           8261,
		PreHandler:     preHandler,
		Routes: []web.Route{
			indexRoute,
			loginRoute,
			loginSubmitRoute,
			logoutRoute,
			signupRoute,
			signupSubmitRoute,
			createKeyRoute,
			createPrivateKeySubmitRoute,
			viewKeyRoute,
			loadKeyRoute,
			importKeyRoute,
			importKeySubmitRoute,
			deleteKeySubmitRoute,
		},
		StaticFilesDir: "web/public",
		TemplatesDir:   "web/templates",
		UseSessions:    true,
	}
	err := ws.Run()
	if err != nil {
		log.Fatal(err)
	}
}
