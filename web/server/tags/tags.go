package tags

import "github.com/jchavannes/jgo/web"

func GetRoutes() []web.Route {
	return []web.Route{
		indexRoute,
		createRoute,
		createSubmitRoute,
	}
}

func preHandler(r *web.Response) {
	r.Helper["Nav"] = "tags"
}
