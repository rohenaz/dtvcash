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

    MemoApp.utf8ByteLength = function (str) {
        // returns the byte length of an utf8 string
        var s = str.length;
        for (var i = str.length - 1; i >= 0; i--) {
            var code = str.charCodeAt(i);
            if (code > 0x7f && code <= 0x7ff) s++;
            else if (code > 0x7ff && code <= 0xffff) s += 2;
            if (code >= 0xDC00 && code <= 0xDFFF) i--; //trail surrogate
        }
        return parseInt(s);
    };

    /**
     * @param {jQuery} $form
     */
    MemoApp.CheckLoadPassword = function($form) {
        if (!localStorage.WalletPassword) {
            return;
        }
        $form.find("[name=password]").val(localStorage.WalletPassword);
        $form.find("[name=save-password]").prop("checked", true);
    };

    /**
     * @param {jQuery} $form
     */
    MemoApp.CheckSavePassword = function($form) {
        if (!$form.find("[name=save-password]").is(':checked')) {
            delete(localStorage.WalletPassword);
            return;
        }

        var password = $form.find("[name=password]").val();
        if (password.length === 0) {
            return;
        }

        localStorage.WalletPassword = password;
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

    /**
     * @param {string} path
     * @param {function} close
     * @return {WebSocket}
     */
    MemoApp.GetSocket = function(path, close) {
        var loc = window.location;
        var protocol = window.location.protocol.toLowerCase() === "https:" ? "wss" : "ws";
        var socket = new WebSocket(protocol + "://" + loc.hostname + ":" + loc.port + path);

        socket.onopen = function () {
            console.log("Socket opened to: " + path);
        };

        socket.onclose = function () {
            console.log("Socket closed to: " + path);
            if (close !== undefined) {
                close();
            }
        };

        setInterval(function () {
            var wsMessage = "heartbeat";
            socket.send(JSON.stringify(wsMessage));
        }, 15000);

        return socket;
    };

    MemoApp.Events = {};

    MemoApp.URL = {
        Index: "",
        Profile: "profile",
        LoadKey: "key/load",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        MemoPost: "post",
        MemoPostAjax: "post-ajax",
        MemoNewSubmit: "memo/new-submit",
        MemoReplySubmit: "memo/reply-submit",
        MemoFollowSubmit: "memo/follow-submit",
        MemoUnfollowSubmit: "memo/unfollow-submit",
        MemoSetNameSubmit: "memo/set-name-submit",
        MemoSetProfileSubmit: "memo/set-profile-submit",
        MemoLikeSubmit: "memo/like-submit",
        MemoWait: "memo/wait",
        MemoWaitSubmit: "memo/wait-submit",
        KeyChangePasswordSubmit: "key/change-password-submit",
        TopicsSocket: "topics/socket",
        TopicsMorePosts: "topics/more-posts",
        TopicsCreateSubmit: "topics/create-submit"
    };
})();
