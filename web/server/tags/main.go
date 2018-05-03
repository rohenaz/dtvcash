package tags

import "github.com/jchavannes/jgo/web"

var urlTagName = web.UrlParam{
	Id:   "tag",
	Type: web.UrlParamAny,
}

func GetRoutes() []web.Route {
	return []web.Route{
		indexRoute,
		createRoute,
		createSubmitRoute,
		viewRoute,
		socketRoute,
	}
}

func preHandler(r *web.Response) {
	r.Helper["Nav"] = "tags"
}
