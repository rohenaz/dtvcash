var MemoApp = {
    Template: {},
    Form: {}
};

(function () {
    /**
     * @param token {string}
     */
    MemoApp.InitCsrf = function (token) {
        /**
         * @param method {string}
         * @returns {boolean}
         */
        function csrfSafeMethod(method) {
            // HTTP methods that do not require CSRF protection.
            return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
        }

        $.ajaxSetup({
            crossDomain: false,
            beforeSend: function (xhr, settings) {
                if (!csrfSafeMethod(settings.type)) {
                    xhr.setRequestHeader("X-CSRF-Token", token);
                }
            }
        });
    };

    MemoApp.InitTimeZone = function () {
        if (document.cookie) {
            var tz = jstz.determine();
            document.cookie = "memo_time_zone=" + tz.name() + ";path=/;max-age=31104000";
        }
    };

    var BaseURL = "/";

    /**
     * @param url {string}
     */
    MemoApp.SetBaseUrl = function (url) {
        BaseURL = url;
    };

    /**
     * @return {string}
     */
    MemoApp.GetBaseUrl = function () {
        return BaseURL;
    };

    /**
     * @param {XMLHttpRequest} xhr
     */
    MemoApp.Form.ErrorHandler = function (xhr) {
        var errorMessage =
            "Error with request (response code " + xhr.status + "):\n" +
            (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
            "If this problem persists, try refreshing the page.";
        alert(errorMessage);
    };

    MemoApp.Events = {};

    MemoApp.URL = {
        Index: "",
        Profile: "profile",
        LoadKey: "key/load",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        MemoPost: "post",
        MemoNewSubmit: "memo/new-submit",
        MemoReplySubmit: "memo/reply-submit",
        MemoFollowSubmit: "memo/follow-submit",
        MemoUnfollowSubmit: "memo/unfollow-submit",
        MemoSetNameSubmit: "memo/set-name-submit",
        MemoSetProfileSubmit: "memo/set-profile-submit",
        MemoLikeSubmit: "memo/like-submit",
        MemoWait: "memo/wait",
        MemoWaitSubmit: "memo/wait-submit",
        KeyChangePasswordSubmit: "key/change-password-submit"
    };
})();
