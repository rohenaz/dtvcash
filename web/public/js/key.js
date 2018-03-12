(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.CreatePrivateKey = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var name = $ele.find("[name=name]").val();
            if (name.length === 0) {
                alert("Must enter a name.");
                return;
            }

            var password = $ele.find("[name=password]").val();
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
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.ImportPrivateKey = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var name = $ele.find("[name=name]").val();
            if (name.length === 0) {
                alert("Must enter a name.");
                return;
            }

            var wif = $ele.find("[name=wif]").val();
            if (wif.length === 0) {
                alert("Must enter a wif.");
                return;
            }

            var password = $ele.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.ImportKeySubmit,
                data: {
                    name: name,
                    wif: wif,
                    password: password
                },
                success: function () {
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Index
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
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
    /**
     * @param {jQuery} $keyLinks
     */
    MemoApp.Form.DeleteKeys = function ($keyLinks) {
        $keyLinks.click(function (e) {
            e.preventDefault();

            var $this = $(this);
            var id = $this.attr("data-id");

            if (!confirm("Are you sure you want to delete this key? (id = " + id + ")")) {
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.DeleteKeySubmit,
                data: {
                    id: id
                },
                success: function () {
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Index
                },
                /**
                 * @param {XMLHttpRequest} xhr
                 */
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
})();
