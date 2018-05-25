package auth

import "github.com/jchavannes/jgo/web"

func GetRoutes() []web.Route {
	return []web.Route{
		loginRoute,
		loginSubmitRoute,
		signupRoute,
		signupSubmitRoute,
		logoutRoute,
	}
}
