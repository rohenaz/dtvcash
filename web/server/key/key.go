package key

import "github.com/jchavannes/jgo/web"

func GetRoutes() []web.Route {
	return []web.Route{
		viewKeyRoute,
		loadKeyRoute,
		changePasswordRoute,
		changePasswordSubmitRoute,
	}
}
