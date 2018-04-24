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
                success: function (txHash) {
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.SetName = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var name = $form.find("[name=name]").val();
            if (name.length === 0) {
                alert("Must enter a name.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoSetNameSubmit,
                data: {
                    name: name,
                    password: password
                },
                success: function (txHash) {
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Follow = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var address = $form.find("[name=address]").val();
            if (address.length === 0) {
                alert("Form error, address not set.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoFollowSubmit,
                data: {
                    address: address,
                    password: password
                },
                success: function (txHash) {
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Unfollow = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var address = $form.find("[name=address]").val();
            if (address.length === 0) {
                alert("Form error, address not set.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoUnfollowSubmit,
                data: {
                    address: address,
                    password: password
                },
                success: function (txHash) {
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Like = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var txHash = $form.find("[name=tx-hash]").val();
            if (txHash.length === 0) {
                alert("Form error, tx hash not set.");
                return;
            }

            var tip = $form.find("[name=tip]").val();
            if (tip.length !== 0 && tip < 546) {
                alert("Must enter a tip greater than 546 (the minimum dust limit).");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoLikeSubmit,
                data: {
                    txHash: txHash,
                    tip: tip,
                    password: password
                },
                success: function (txHash) {
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Wait = function ($form) {
        $form.submit(function (e) {
            e.preventDefault();
            var txHash = $form.find("[name=tx-hash]").val();
            if (txHash.length === 0) {
                alert("Form error, tx hash not set.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                data: {
                    txHash: txHash
                },
                success: function (url) {
                    if (url === undefined || url.length === 0) {
                        alert("Error with broadcast. Please try again.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + url
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
        $form.submit();
    };
})();
