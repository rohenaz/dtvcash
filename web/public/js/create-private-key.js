(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.CreatePrivateKey = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var name = $ele.find("[name=name]").val();
            var password = $ele.find("[name=password]").val();

            if (name.length === 0) {
                alert("Must enter a name.");
                return;
            }

            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.CreatePrivateKeySubmit,
                data: {
                    name: name,
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
                        "Error creating private key (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
})();
