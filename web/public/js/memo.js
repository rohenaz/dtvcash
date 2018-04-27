(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.LogoutButton = function ($ele) {
        $ele.click(function() {
            delete(localStorage.WalletPassword);
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewMemo = function ($form) {
        CheckLoadPassword($form);
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

            CheckSavePassword($form);

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoNewSubmit,
                data: {
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
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
        CheckLoadPassword($form);
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

            CheckSavePassword($form);

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoSetNameSubmit,
                data: {
                    name: name,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
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
        CheckLoadPassword($form);
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

            CheckSavePassword($form);

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoFollowSubmit,
                data: {
                    address: address,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
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
        CheckLoadPassword($form);
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

            CheckSavePassword($form);

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
        CheckLoadPassword($form);
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

            CheckSavePassword($form);

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoLikeSubmit,
                data: {
                    txHash: txHash,
                    tip: tip,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
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
    MemoApp.Form.ReplyMemo = function ($form) {
        CheckLoadPassword($form);
        $form.submit(function (e) {
            e.preventDefault();
            var txHash = $form.find("[name=tx-hash]").val();
            if (txHash.length === 0) {
                alert("Form error, tx hash not set.");
                return;
            }

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
            CheckSavePassword($form);

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoReplySubmit,
                data: {
                    txHash: txHash,
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: MemoApp.Form.ErrorHandler
            });
        });
    };
    function CheckLoadPassword($form) {
        if (!localStorage.WalletPassword) {
            return;
        }
        $form.find("[name=password]").val(localStorage.WalletPassword);
        $form.find("[name=save-password]").prop("checked", true);
    }

    /**
     * @param {jQuery} $form
     */
    function CheckSavePassword($form) {
        if (!$form.find("[name=save-password]").is(':checked')) {
            delete(localStorage.WalletPassword);
            return;
        }

        var password = $form.find("[name=password]").val();
        if (password.length === 0) {
            return;
        }

        localStorage.WalletPassword = password;
    }
    /**
     * @param {jQuery} $form
     * @param {jQuery} $notify
     * @param {jQuery} $title
     */
    MemoApp.Form.Wait = function ($form, $notify, $title) {
        var text = "Broadcasting transaction";
        var dots = 1;
        setInterval(function() {
            $title.html(text + Array(dots).join("."));
            dots++;
            if (dots > 5) {
                dots = 1;
            }
        }, 750);
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
                    if (!url || url.length === 0) {
                        alert("Error with broadcast. Please try again.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + url
                },
                error: function () {
                    $notify.html(
                        "Transaction propagation taking longer than normal. " +
                        "You can continue waiting or try again. " +
                        "This page will automatically redirect when transaction has propagated."
                    );
                    $form.submit();
                }
            });
        });
        $form.submit();
    };

    /**
     * @param {jQuery} $likesButton
     * @param {jQuery} $likes
     */
    MemoApp.Form.LikesToggle = function($likesButton, $likes) {
        $likesButton.click(function (e) {
            e.preventDefault();
            if ($likes.is(":visible")) {
                $likes.hide();
                $likesButton.html("Show");
            } else {
                $likes.show();
                $likesButton.html("Hide");
            }
        });
    };

    MemoApp.utf8ByteLength = function(str) {
        // returns the byte length of an utf8 string
        var s = str.length;
        for (var i=str.length-1; i>=0; i--) {
            var code = str.charCodeAt(i);
            if (code > 0x7f && code <= 0x7ff) s++;
            else if (code > 0x7ff && code <= 0xffff) s+=2;
            if (code >= 0xDC00 && code <= 0xDFFF) i--; //trail surrogate
        }
        return parseInt(s);
    }
})();
