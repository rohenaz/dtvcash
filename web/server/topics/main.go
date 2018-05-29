package topics

import "github.com/jchavannes/jgo/web"

var urlTopicName = web.UrlParam{
	Id:   "topic",
	Type: web.UrlParamAny,
}

var urlTxHash = web.UrlParam{
	Id:   "tx-hash",
	Type: web.UrlParamString,
}

func GetRoutes() []web.Route {
	return []web.Route{
		indexRoute,
		createRoute,
		createSubmitRoute,
		viewRoute,
		socketRoute,
		postsMoreRoute,
		postAjaxRoute,
	}
}

func preHandler(r *web.Response) {
	r.Helper["Nav"] = "channels"
}
