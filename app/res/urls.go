package res

import "github.com/jchavannes/jgo/web"

const (
	UrlIndex        = "/"
	UrlSignup       = "/signup"
	UrlSignupSubmit = "/signup-submit"
	UrlLogin        = "/login"
	UrlLoginSubmit  = "/login-submit"
	UrlLogout       = "/logout"
	UrlTests        = "/tests"
	UrlProtocol     = "/protocol"
	UrlDisclaimer   = "/disclaimer"
)

const (
	UrlKeyExport = "/key/export"
	UrlKeyLoad   = "/key/load"
)

const (
	UrlMemoNew           = "/memo/new"
	UrlMemoNewSubmit     = "/memo/new-submit"
	UrlMemoSetName       = "/memo/set-name"
	UrlMemoSetNameSubmit = "/memo/set-name-submit"
	UrlMemoFollow        = "/memo/follow"
	UrlMemoFollowSumbit  = "/memo/follow-submit"
)

const (
	UrlProfiles    = "/profiles"
	UrlProfileView = "/profile"

	TmplProfiles = "/profile/all"
)

func GetBaseUrl(r *web.Response) string {
	baseUrl := r.Request.GetHeader("AppPath")
	if baseUrl == "" {
		baseUrl = "/"
	}
	return baseUrl
}

func GetUrlWithBaseUrl(url string, r *web.Response) string {
	baseUrl := GetBaseUrl(r)
	baseUrl = baseUrl[:len(baseUrl)-1]
	return baseUrl + url
}
