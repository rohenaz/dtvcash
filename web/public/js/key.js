(function () {
    /**
     * @param {int} id
     * @param {jQuery} $form
     * @param {jQuery} $keyDiv
     */
    MemoApp.Form.LoadKey = function (id, $form, $keyDiv) {
        $form.submit(function (e) {
            e.preventDefault();
            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.LoadKey,
                data: {
                    id: id,
                    password: password
                },
                success: function (keyHtml) {
                    $keyDiv.html(keyHtml);
                },
                /**
                 * @param {XMLHttpRequest} xhr
                 */
                error: function (xhr) {
                    if (xhr.status === 401) {
                        alert("Error unlocking. Please try again.");
                    } else {
                        MemoApp.Form.ErrorHandler(xhr);
                    }
                }
            });
        });
    };
})();
