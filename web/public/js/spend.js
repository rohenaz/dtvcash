(function () {
    /**
     * @param {int} utxoId
     * @param {jQuery} $form
     * @param {jQuery} $outDiv
     */
    MemoApp.Form.Spend = function (utxoId, $form, $outDiv) {
        $form.submit(function (e) {
            e.preventDefault();
            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.SpendSign,
                data: {
                    id: utxoId,
                    password: password
                },
                success: function (keyHtml) {
                    $outDiv.html(keyHtml);
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
