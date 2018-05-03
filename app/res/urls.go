package res

import "github.com/jchavannes/jgo/web"

const (
	UrlIndex        = "/"
	UrlSignup       = "/signup"
	UrlSignupSubmit = "/signup-submit"
	UrlLogin        = "/login"
	UrlLoginSubmit  = "/login-submit"
	UrlLogout       = "/logout"
	UrlProtocol     = "/protocol"
	UrlDisclaimer   = "/disclaimer"
	UrlIntroducing  = "/introducing-memo"
	UrlNeedFunds    = "/need-funds"
	UrlNewPosts     = "/new-posts"
	UrlStats        = "/stats"
	UrlFeed         = "/feed"
	UrlAbout        = "/about"
)

const (
	UrlKeyExport               = "/key/export"
	UrlKeyLoad                 = "/key/load"
	UrlKeyChangePassword       = "/key/change-password"
	UrlKeyChangePasswordSubmit = "/key/change-password-submit"
)

const (
	UrlMemoNew              = "/memo/new"
	UrlMemoNewSubmit        = "/memo/new-submit"
	UrlMemoSetName          = "/memo/set-name"
	UrlMemoSetNameSubmit    = "/memo/set-name-submit"
	UrlMemoFollow           = "/memo/follow"
	UrlMemoFollowSubmit     = "/memo/follow-submit"
	UrlMemoUnfollow         = "/memo/unfollow"
	UrlMemoUnfollowSubmit   = "/memo/unfollow-submit"
	UrlMemoPost             = "/post"
	UrlMemoLike             = "/memo/like"
	UrlMemoLikeSubmit       = "/memo/like-submit"
	UrlMemoReply            = "/memo/reply"
	UrlMemoReplySubmit      = "/memo/reply-submit"
	UrlMemoWait             = "/memo/wait"
	UrlMemoWaitSubmit       = "/memo/wait-submit"
	UrlMemoSetProfile       = "/memo/set-profile"
	UrlMemoSetProfileSubmit = "/memo/set-profile-submit"

	TmplMemoPost = "/memo/post"
)

const (
	UrlProfiles         = "/profiles"
	UrlProfileView      = "/profile"
	UrlProfileFollowers = "/profile/followers"
	UrlProfileFollowing = "/profile/following"

	TmplProfiles = "/profile/all"
)

const (
	UrlPostsNew          = "/posts/new"
	UrlPostsTop          = "/posts/top"
	UrlPostsArchive      = "/posts/archive"
	UrlPostsPersonalized = "/posts/personalized"
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
