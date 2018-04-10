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
)

const (
	UrlKeyView         = "/key"
	UrlKeyLoad         = "/key/load"
	UrlKeyImport       = "/key/import"
	UrlKeyImportSubmit = "/key/import-submit"
	UrlKeyCreate       = "/key/create"
	UrlKeyCreateSubmit = "/key/create-submit"
	UrlKeyDeleteSubmit = "/key/delete-submit"
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
