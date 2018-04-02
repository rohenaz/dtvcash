package key

import "github.com/jchavannes/jgo/web"

var urlId = web.UrlParam{
	Id:   "id",
	Type: web.UrlParamInteger,
}

func GetRoutes() []web.Route {
	return []web.Route{
		createKeyRoute,
		createPrivateKeySubmitRoute,
		viewKeyRoute,
		loadKeyRoute,
		importKeyRoute,
		importKeySubmitRoute,
		deleteKeySubmitRoute,
		dataLoadSubmitRoute,
		refreshKeySubmitRoute,
	}
}
