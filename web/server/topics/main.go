package topics

import "github.com/jchavannes/jgo/web"

var urlTopicName = web.UrlParam{
	Id:   "topic",
	Type: web.UrlParamAny,
}

func GetRoutes() []web.Route {
	return []web.Route{
		indexRoute,
		createRoute,
		createSubmitRoute,
		viewRoute,
		socketRoute,
		postsMoreRoute,
	}
}

func preHandler(r *web.Response) {
	r.Helper["Nav"] = "topics"
}
