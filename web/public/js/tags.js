(function () {

    var maxNewTagBytes = 74;

    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewTag = function ($form) {
        var $tagName = $form.find("[name=tag]");
        var $message = $form.find("[name=message]");

        var $tagByteCount = $form.find(".name-byte-count");
        var $msgByteCount = $form.find(".message-byte-count");

        $tagName.on("input", function () {
            setMsgByteCount();
        });
        $message.on("input", function () {
            setMsgByteCount();
        });

        function getByteSize() {
            return MemoApp.utf8ByteLength($tagName.val()) + MemoApp.utf8ByteLength($message.val());
        }

        function setMsgByteCount() {
            var cnt = maxNewTagBytes - getByteSize();
            $tagByteCount.html("[" + cnt + "]");
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $tagByteCount.addClass("red");
                $msgByteCount.addClass("red");
            } else {
                $tagByteCount.removeClass("red");
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

            var tagName = $tagName.val();
            var message = $message.val();
            if (maxNewTagBytes - getByteSize() < 0) {
                alert("Maximum size is " + maxNewTagBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (tagName.length === 0) {
                alert("Must enter a tag name.");
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
                url: MemoApp.GetBaseUrl() + MemoApp.URL.TagsCreateSubmit,
                data: {
                    tag: tagName,
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
     * @param {string} tag
     * @param {jQuery} $morePosts
     */
    MemoApp.WatchNewTags = function (tag, $morePosts) {
        socket = MemoApp.GetSocket(MemoApp.GetBaseUrl() + MemoApp.URL.TagsSocket + "?tag=" + tag, function () {
            console.log("socket closed...");
        });
        /**
         * @param {MessageEvent} msg
         */
        socket.onmessage = function (msg) {
            console.log(msg.data);
            var txHash = msg.data.replace(/['"]+/g, '')
            $.ajax({
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoPostAjax + "/" + txHash,
                success: function(html) {
                    $morePosts.append(html);
                },
                error: function (xhr) {
                    alert("error getting post via ajax (status: " + xhr.status + ")");
                }
            });
        };
    };

    /**
     * @param {jQuery} $broadcasting
     * @param {jQuery} $form
     */
    MemoApp.Form.NewTagMessage = function ($broadcasting, $form) {
        var $tagName = $form.find("[name=tag]");
        var $message = $form.find("[name=message]");

        var $msgByteCount = $form.find(".message-byte-count");

        $message.on("input", function () {
            setMsgByteCount();
        });

        function getByteSize() {
            return MemoApp.utf8ByteLength($tagName.val()) + MemoApp.utf8ByteLength($message.val());
        }

        function setMsgByteCount() {
            var cnt = maxNewTagBytes - getByteSize();
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

            var tagName = $tagName.val();
            var message = $message.val();
            if (maxNewTagBytes - getByteSize() < 0) {
                alert("Maximum size is " + maxNewTagBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (tagName.length === 0) {
                alert("Must enter a tag name.");
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
                url: MemoApp.GetBaseUrl() + MemoApp.URL.TagsCreateSubmit,
                data: {
                    tag: tagName,
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    submitting = false;
                    if (!txHash || txHash.length === 0) {
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    $broadcasting.show();
                    $form.hide();
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: txHash
                        },
                        success: function () {
                            $broadcasting.hide();
                            $message.val("");
                            setMsgByteCount();
                            $form.show();
                        },
                        error: function () {
                            alert("Error waiting for transaction to broadcast.");
                            $broadcasting.hide();
                            $message.val("");
                            $form.show();
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
})();
