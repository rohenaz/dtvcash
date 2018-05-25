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
	UrlGuides       = "/guides"
	UrlDisclaimer   = "/disclaimer"
	UrlIntroducing  = "/introducing-memo"
	UrlOpenSource   = "/open-sourcing-memo"
	UrlNeedFunds    = "/need-funds"
	UrlNewPosts     = "/new-posts"
	UrlStats        = "/stats"
	UrlFeed         = "/feed"
	UrlAbout        = "/about"
	UrlNotFound     = "/404"
)

const (
	UrlKeyExport               = "/key/export"
	UrlKeyLoad                 = "/key/load"
	UrlKeyChangePassword       = "/key/change-password"
	UrlKeyChangePasswordSubmit = "/key/change-password-submit"
	UrlKeyDeleteAccount        = "/key/delete-account"
	UrlKeyDeleteAccountSubmit  = "/key/delete-account-submit"
)

const (
	UrlMemoNew                  = "/memo/new"
	UrlMemoNewSubmit            = "/memo/new-submit"
	UrlMemoSetName              = "/memo/set-name"
	UrlMemoSetNameSubmit        = "/memo/set-name-submit"
	UrlMemoFollow               = "/memo/follow"
	UrlMemoFollowSubmit         = "/memo/follow-submit"
	UrlMemoUnfollow             = "/memo/unfollow"
	UrlMemoUnfollowSubmit       = "/memo/unfollow-submit"
	UrlMemoPost                 = "/post"
	UrlMemoPostThreaded         = "/post-threaded"
	UrlMemoPostMoreThreadedAjax = "/post-more-threaded-ajax"
	UrlMemoPostThreadedAjax     = "/post-threaded-ajax"
	UrlMemoPostAjax             = "/post-ajax"
	UrlMemoLike                 = "/memo/like"
	UrlMemoLikeSubmit           = "/memo/like-submit"
	UrlMemoReply                = "/memo/reply"
	UrlMemoReplySubmit          = "/memo/reply-submit"
	UrlMemoWait                 = "/memo/wait"
	UrlMemoWaitSubmit           = "/memo/wait-submit"
	UrlMemoSetProfile           = "/memo/set-profile"
	UrlMemoSetProfileSubmit     = "/memo/set-profile-submit"

	TmplMemoPost         = "/memo/post"
	TmplMemoPostThreaded = "/memo/post-threaded"
)

const (
	UrlProfiles              = "/profiles"
	UrlProfileView           = "/profile"
	UrlProfileFollowers      = "/profile/followers"
	UrlProfileFollowing      = "/profile/following"
	UrlProfileSettings       = "/settings"
	UrlProfileSettingsSubmit = "/settings-submit"
	UrlProfileNotifications  = "/notifications"

	TmplProfiles             = "/profile/all"
	TmplProfileSettings      = "/profile/settings"
	TmplProfileNotifications = "/profile/notifications"
)

const (
	UrlPostsNew          = "/posts/new"
	UrlPostsTop          = "/posts/top"
	UrlPostsRanked       = "/posts/ranked"
	UrlPostsArchive      = "/posts/archive"
	UrlPostsPersonalized = "/posts/personalized"
)

const (
	UrlTopics             = "/topics"
	UrlTopicsCreate       = "/topics/create"
	UrlTopicsCreateSubmit = "/topics/create-submit"
	UrlTopicView          = "/topic"
	UrlTopicsSocket       = "/topics/socket"
	UrlTopicsMorePosts    = "/topics/more-posts"
	UrlTopicsPostAjax     = "/topics/post-ajax"

	TmplTopicView  = "/topics/view"
	TmplTopicPost  = "/topics/post"
	TmplTopicPosts = "/topics/posts"
)

const (
	TmplSnippetsPost                 = "/snippets/post"
	TmplSnippetsPostThreaded         = "/snippets/post-threaded"
	TmplSnippetsPostThreadedLoadMore = "/snippets/post-threaded-load-more"
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
