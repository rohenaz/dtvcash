(function () {

    var maxPostBytes = 217;
    var maxReplyBytes = 184;
    var maxNameBytes = 77;
    var maxProfileTextBytes = 217;

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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
     * @param {string} txHash
     * @param {boolean} threaded
     */
    MemoApp.Form.ReplyMemo = function (txHash, threaded) {
        var $post = $("#post-" + txHash);
        var $form = $("#reply-form-" + txHash);
        var $replyCancel = $("#reply-cancel-" + txHash);
        var $message = $form.find("[name=message]");
        var $msgByteCount = $form.find(".message-byte-count");
        var $replyLink = $("#reply-link-" + txHash);
        var $broadcasting = $post.find(".broadcasting:eq(0)");
        var $creating = $post.find(".creating:eq(0)");
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

        $replyCancel.click(function(e) {
            e.preventDefault();
            $form.addClass("hidden");
        });

        setMsgByteCount();
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

            $creating.removeClass("hidden");
            $replyLink.hide();
            $form.hide();

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoReplySubmit,
                data: {
                    txHash: txHash,
                    message: message,
                    password: password
                },
                success: function (replyTxHash) {
                    submitting = false;
                    if (!replyTxHash || replyTxHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    $creating.addClass("hidden");
                    $broadcasting.removeClass("hidden");
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: replyTxHash
                        },
                        success: function () {
                            submitting = false;
                            var url = MemoApp.URL.MemoPostAjax;
                            if (threaded) {
                                url = MemoApp.URL.MemoPostThreadedAjax
                            }
                            $.ajax({
                                url: MemoApp.GetBaseUrl() + url + "/" + txHash,
                                success: function (html) {
                                    $("#post-" + txHash).replaceWith(html);
                                },
                                error: function (xhr) {
                                    alert("error getting post via ajax (status: " + xhr.status + ")");
                                }
                            });
                        },
                        error: function () {
                            submitting = false;
                            $broadcasting.addClass("hidden");
                            console.log("Error waiting for transaction to broadcast.");
                        }
                    });
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
     * @param {string} txHash
     */
    MemoApp.Form.ReplyLink = function(txHash) {
        var $replyLink = $("#reply-link-" + txHash);
        var $replyForm = $("#reply-form-" + txHash);
        $replyLink.click(function(e) {
            e.preventDefault();
            $replyForm.removeClass("hidden");
        });
    };

    /**
     * @param {jQuery} $like
     * @param {string} txHash
     * @param {boolean} threaded
     */
    MemoApp.Form.NewLike = function ($like, txHash, threaded) {
        var $likeLink = $("#like-link-" + txHash);
        var $likeCancel = $("#like-cancel-" + txHash);
        var $likeInfo = $("#like-info-" + txHash);
        var $likeForm = $("#like-form-" + txHash);
        var $creating = $like.parent().find(".creating:eq(0)");
        var $broadcasting = $like.parent().find(".broadcasting:eq(0)");

        $likeLink.click(function (e) {
            e.preventDefault();
            $likeInfo.hide();
            $likeForm.removeClass("hidden");
        });
        $likeCancel.click(function (e) {
            e.preventDefault();
            $likeInfo.show();
            $likeForm.addClass("hidden");
        });

        var submitting = false;
        $likeForm.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var tip = $likeForm.find("[name=tip]").val();
            if (tip.length !== 0 && tip < 546) {
                alert("Must enter a tip greater than 546 (the minimum dust limit).");
                return;
            }
            $creating.removeClass("hidden");
            $likeForm.hide();

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoLikeSubmit,
                data: {
                    txHash: txHash,
                    tip: tip,
                    password: password
                },
                success: function (likeTxHash) {
                    submitting = false;
                    if (!likeTxHash || likeTxHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    $creating.addClass("hidden");
                    $broadcasting.removeClass("hidden");
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: likeTxHash
                        },
                        success: function () {
                            submitting = false;
                            var url = MemoApp.URL.MemoPostAjax;
                            if (threaded) {
                                url = MemoApp.URL.MemoPostThreadedAjax
                            }
                            $.ajax({
                                url: MemoApp.GetBaseUrl() + url + "/" + txHash,
                                success: function (html) {
                                    $("#post-" + txHash).replaceWith(html);
                                },
                                error: function (xhr) {
                                    alert("error getting post via ajax (status: " + xhr.status + ")");
                                }
                            });
                        },
                        error: function () {
                            submitting = false;
                            $broadcasting.addClass("hidden");
                            console.log("Error waiting for transaction to broadcast.");
                        }
                    });
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

    /**
     * @param {jQuery} $moreReplies
     * @param {string} txHash
     * @param {number} offset
     */
    MemoApp.Form.LoadMoreReplies = function ($moreReplies, txHash, offset) {
        var $link = $moreReplies.find("a");
        $link.click(function (e) {
            e.preventDefault();
            $.ajax({
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoPostMoreThreadedAjax,
                data: {
                    txHash: txHash,
                    offset: offset + 25
                },
                success: function (html) {
                    $moreReplies.replaceWith(html);
                },
                error: function () {
                    console.log("Error loading more replies.");
                }
            });
        });
    };
})();
