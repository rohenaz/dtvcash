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

    MemoApp.Events = {};

    MemoApp.URL = {
        Index: "",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        CreatePrivateKeySubmit: "create-private-key-submit"
    };
})();
