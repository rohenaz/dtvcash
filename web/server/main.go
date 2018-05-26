package server

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/queuer"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	auth2 "github.com/memocash/memo/web/server/auth"
	"github.com/memocash/memo/web/server/key"
	"github.com/memocash/memo/web/server/memo"
	"github.com/memocash/memo/web/server/poll"
	"github.com/memocash/memo/web/server/posts"
	"github.com/memocash/memo/web/server/profile"
	"github.com/memocash/memo/web/server/topics"
	"github.com/nicksnyder/go-i18n/i18n"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var UseMinJS bool

func isLoggedIn(r *web.Response) bool {
	if ! auth.IsLoggedIn(r.Session.CookieId) {
		r.SetRedirect(res.UrlLogin)
		return false
	}
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
		userAddress, err := cache.GetUserAddress(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting user address from cache", err), http.StatusInternalServerError)
			return
		}
		r.Helper["UserAddress"] = userAddress.GetEncoded()
		userSettings, err := cache.GetUserSettings(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting user settings from cache", err), http.StatusInternalServerError)
			return
		}
		unreadNotifications, err := cache.GetUnreadNotificationCount(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting last notification id from cache", err), http.StatusInternalServerError)
			return
		}
		r.Helper["UnreadNotifications"] = unreadNotifications
		r.Helper["UserSettings"] = userSettings
	} else {
		r.Helper["UserSettings"] = db.GetDefaultUserSettings()
	}
	if UseMinJS {
		r.Helper["jsFiles"] = res.GetMinJsFiles()
	} else {
		r.Helper["jsFiles"] = res.GetResJsFiles()
	}
	r.Helper["cssFiles"] = res.GetResCssFiles()
	r.Helper["TimeZone"] = r.Request.GetCookie("memo_time_zone")
	r.Helper["Nav"] = ""

	lang := r.Request.GetCookie("memo_language")
	if lang == "" {
		lang = r.Request.GetHeader("Accept-Language")
	}
	if ! isValidLang(lang) {
		lang = "en-US"
	}

	r.SetFuncMap(map[string]interface{}{
		"T": i18n.MustTfunc(lang),
		"Title": strings.Title,
	})
}

func notFoundHandler(r *web.Response) {
	r.SetResponseCode(http.StatusNotFound)
	r.RenderTemplate(res.UrlNotFound)
}

func isValidLang(lang string) bool {
	for _, item := range []string{"en-US", "es-LA", "zh-CN", "ja-JP", "fr-FR"} {
		if item == lang {
			return true
		}
	}
	return false
}

var allowedExtensions = []string{
	"js",
	"css",
	"jpg",
	"png",
	"ico",
	"gif",
	"woff",
	"woff2",
	"ttf",
	"svg",
	"eot",
}

func Run(sessionCookieInsecure bool) {
	go func() {
		queuer.StartAndKeepAlive()
	}()

	var langDir = "web/lang"
	files, err := ioutil.ReadDir(langDir)
	if err != nil {
		log.Fatal(jerr.Get("error getting language files", err))
	}

	for _, file := range files {
		i18n.MustLoadTranslationFile(langDir + "/" + file.Name())
	}

	// Start web server
	ws := web.Server{
		CookiePrefix:      "memo",
		InsecureCookie:    sessionCookieInsecure,
		AllowedExtensions: allowedExtensions,
		IsLoggedIn:        isLoggedIn,
		Port:              8261,
		NotFoundHandler:   notFoundHandler,
		PreHandler:        preHandler,
		GetCsrfToken:      getCsrfToken,
		Routes: web.Routes(
			[]web.Route{
				indexRoute,
				protocolRoute,
				guidesRoute,
				disclaimerRoute,
				introducingMemoRoute,
				openSourcingMemoRoute,
				aboutRoute,
				needFundsRoute,
				newPostsRoute,
				statsRoute,
				feedRoute,
				//testsRoute,
			},
			poll.GetRoutes(),
			topics.GetRoutes(),
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
	err = ws.Run()
	if err != nil {
		log.Fatal(err)
	}
}
