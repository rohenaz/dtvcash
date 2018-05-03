package server

import (
	"fmt"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/queuer"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	auth2 "git.jasonc.me/main/memo/web/server/auth"
	"git.jasonc.me/main/memo/web/server/key"
	"git.jasonc.me/main/memo/web/server/memo"
	"git.jasonc.me/main/memo/web/server/posts"
	"git.jasonc.me/main/memo/web/server/profile"
	"git.jasonc.me/main/memo/web/server/tags"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"log"
	"net/http"
)

var UseMinJS bool

func isLoggedIn(r *web.Response) bool {
	if ! auth.IsLoggedIn(r.Session.CookieId) {
		fmt.Println("here1")
		r.SetRedirect(res.UrlLogin)
		return false
	}
	fmt.Println("here2")
	return true
}

func getCsrfToken(cookieId string) string {
	token, err := db.GetCsrfTokenString(cookieId)
	if err != nil {
		jerr.Get("error getting csrf token", err).Print()
		return ""
	}
	return token
}

var blockedIps = []string{
	"91.130.64.132",
	"49.195.117.8",
}

func preHandler(r *web.Response) {
	for _, blockedIp := range blockedIps {
		if r.Request.GetSourceIP() == blockedIp {
			r.Error(jerr.Newf("blocked ip: %s\n", blockedIp), http.StatusUnauthorized)
			return
		}
	}
	r.Helper["Title"] = "Memo"
	r.Helper["Description"] = "Decentralized on-chain social network built on Bitcoin Cash"
	r.Helper["BaseUrl"] = res.GetBaseUrl(r)
	if r.Request.HttpRequest.Host != "memo.cash" {
		r.Helper["Dev"] = true
		r.Helper["GoogleId"] = "UA-23518512-10"
	} else {
		r.Helper["Dev"] = false
		r.Helper["GoogleId"] = "UA-23518512-9"
	}
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
	r.Helper["cssFiles"] = res.GetResCssFiles()
	r.Helper["TimeZone"] = r.Request.GetCookie("memo_time_zone")
	r.Helper["Nav"] = ""
}

func Run(sessionCookieInsecure bool) {
	go func() {
		queuer.StartAndKeepAlive()
	}()

	// Start web server
	ws := web.Server{
		CookiePrefix:   "memo",
		InsecureCookie: sessionCookieInsecure,
		IsLoggedIn:     isLoggedIn,
		Port:           8261,
		PreHandler:     preHandler,
		GetCsrfToken:   getCsrfToken,
		Routes: web.Routes(
			[]web.Route{
				indexRoute,
				protocolRoute,
				disclaimerRoute,
				introducingMemoRoute,
				aboutRoute,
				needFundsRoute,
				newPostsRoute,
				statsRoute,
				feedRoute,
				//testsRoute,
			},
			tags.GetRoutes(),
			posts.GetRoutes(),
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
