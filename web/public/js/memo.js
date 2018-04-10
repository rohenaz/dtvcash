(function () {
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewMemo = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var message = $form.find("[name=message]").val();
            if (message.length === 0) {
                alert("Must enter a message.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoNewSubmit,
                data: {
                    message: message,
                    password: password
                },
                success: function () {
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Index
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
})();
