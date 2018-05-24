package poll

import "github.com/jchavannes/jgo/web"

func GetRoutes() []web.Route {
	return []web.Route{
		createRoute,
		createSubmitRoute,
		voteSubmitRoute,
	}
}

