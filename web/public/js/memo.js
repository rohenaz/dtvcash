(function () {

    var maxPostBytes = 77;
    var maxReplyBytes = 45;
    var maxNameBytes = 77;
    var maxProfileTextBytes = 77;

    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.LogoutButton = function ($ele) {
        $ele.click(function () {
            delete(localStorage.WalletPassword);
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewMemo = function ($form) {
        var $message = $form.find("[name=message]");
        var $msgByteCount = $form.find(".message-byte-count");
        $message.on("input", function () {
            setMsgByteCount();
        });

        function setMsgByteCount() {
            var cnt = maxPostBytes - MemoApp.utf8ByteLength($message.val());
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $msgByteCount.addClass("red");
            } else {
                $msgByteCount.removeClass("red");
            }
        }

        setMsgByteCount();
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var message = $message.val();
            if (maxPostBytes - MemoApp.utf8ByteLength(message) < 0) {
                alert("Maximum post message is " + maxPostBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (message.length === 0) {
                alert("Must enter a message.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoNewSubmit,
                data: {
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.SetName = function ($form) {
        var $name = $form.find("[name=name]");
        var $msgByteCount = $form.find(".message-byte-count");
        $name.on("input", function () {
            setMsgByteCount();
        });

        function setMsgByteCount() {
            var cnt = maxNameBytes - MemoApp.utf8ByteLength($name.val());
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $msgByteCount.addClass("red");
            } else {
                $msgByteCount.removeClass("red");
            }
        }

        setMsgByteCount();
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var name = $name.val();
            if (maxNameBytes - MemoApp.utf8ByteLength(name) < 0) {
                alert("Maximum name is " + maxNameBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (name.length === 0) {
                alert("Must enter a name.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoSetNameSubmit,
                data: {
                    name: name,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.SetProfile = function ($form) {
        var $profile = $form.find("[name=profile]");
        var $msgByteCount = $form.find(".message-byte-count");
        $profile.on("input", function () {
            setMsgByteCount();
        });

        function setMsgByteCount() {
            var cnt = maxProfileTextBytes - MemoApp.utf8ByteLength($profile.val());
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $msgByteCount.addClass("red");
            } else {
                $msgByteCount.removeClass("red");
            }
        }

        setMsgByteCount();

        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var profile = $profile.val();
            if (maxProfileTextBytes - MemoApp.utf8ByteLength(profile) < 0) {
                alert("Maximum profile text is " + maxProfileTextBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (profile.length === 0) {
                if (!confirm("Are you sure you want to set an empty profile?")) {
                    return;
                }
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoSetProfileSubmit,
                data: {
                    profile: profile,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Follow = function ($form) {
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

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

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoFollowSubmit,
                data: {
                    address: address,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };

    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Unfollow = function ($form) {
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

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

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoUnfollowSubmit,
                data: {
                    address: address,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (txHash === undefined || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Like = function ($form) {
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

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

            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoLikeSubmit,
                data: {
                    txHash: txHash,
                    tip: tip,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.ReplyMemo = function ($form) {
        var $message = $form.find("[name=message]");
        var $msgByteCount = $form.find(".message-byte-count");
        $message.on("input", function () {
            setMsgByteCount();
        });

        function setMsgByteCount() {
            var cnt = maxReplyBytes - MemoApp.utf8ByteLength($message.val());
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $msgByteCount.addClass("red");
            } else {
                $msgByteCount.removeClass("red");
            }
        }

        setMsgByteCount();
        MemoApp.CheckLoadPassword($form);
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var message = $message.val();
            if (maxReplyBytes - MemoApp.utf8ByteLength(message) < 0) {
                alert("Maximum reply message is " + maxReplyBytes + " bytes. Note that some characters are more than 1 byte. " +
                    "Emojis are usually 4 bytes, for example.");
                return;
            }

            if (message.length === 0) {
                alert("Must enter a message.");
                return;
            }

            var txHash = $form.find("[name=tx-hash]").val();
            if (txHash.length === 0) {
                alert("Form error, tx hash not set.");
                return;
            }

            var password = $form.find("[name=password]").val();
            if (password.length === 0) {
                alert("Must enter a password.");
                return;
            }
            MemoApp.CheckSavePassword($form);

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoReplySubmit,
                data: {
                    txHash: txHash,
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.MemoWait + "/" + txHash
                },
                error: function (xhr) {
                    submitting = false;
                    if (xhr.status === 401) {
                        alert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    alert(errorMessage);
                }
            });
        });
    };

    /**
     * @param {jQuery} $form
     * @param {jQuery} $notify
     * @param {jQuery} $title
     */
    MemoApp.Form.Wait = function ($form, $notify, $title) {
        var text = "Broadcasting transaction";
        var dots = 1;
        setInterval(function () {
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
    MemoApp.Form.LikesToggle = function ($likesButton, $likes) {
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
})();
