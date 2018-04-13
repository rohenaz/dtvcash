package memo

import "github.com/jchavannes/jgo/web"

var urlAddress = web.UrlParam{
	Id:   "address",
	Type: web.UrlParamString,
}

func GetRoutes() []web.Route {
	return []web.Route{
		newRoute,
		newSubmitRoute,
		setNameRoute,
		setNameSubmitRoute,
		followRoute,
		followSubmitRoute,
		unfollowRoute,
		unfollowSubmitRoute,
	}
}
