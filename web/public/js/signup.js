(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.Signup = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var username = $ele.find("[name=username]").val();
            var password = $ele.find("[name=password]").val();

            if (username.length === 0) {
                alert("Must enter a username.");
                return;
            }

            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.SignupSubmit,
                data: {
                    username: username,
                    password: password
                },
                success: function () {
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Index
                },
                /**
                 * @param {XMLHttpRequest} xhr
                 */
                error: function (xhr) {
                    var errorMessage =
                        "Error creating account:\n" + xhr.responseText + "\n" +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
})();
