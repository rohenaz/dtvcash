(function () {

    var maxNewTopicBytes = 74;

    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewTopic = function ($form) {
        var $topicName = $form.find("[name=topic]");
        var $message = $form.find("[name=message]");

        var $topicByteCount = $form.find(".name-byte-count");
        var $msgByteCount = $form.find(".message-byte-count");

        $topicName.on("input", function () {
            setMsgByteCount();
        });
        $message.on("input", function () {
            setMsgByteCount();
        });

        function getByteSize() {
            return MemoApp.utf8ByteLength($topicName.val()) + MemoApp.utf8ByteLength($message.val());
        }

        function setMsgByteCount() {
            var cnt = maxNewTopicBytes - getByteSize();
            $topicByteCount.html("[" + cnt + "]");
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $topicByteCount.addClass("red");
                $msgByteCount.addClass("red");
            } else {
                $topicByteCount.removeClass("red");
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

            var topicName = $topicName.val();
            var message = $message.val();
            if (maxNewTopicBytes - getByteSize() < 0) {
                alert("Maximum size is " + maxNewTopicBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (topicName.length === 0) {
                alert("Must enter a topic name.");
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
                url: MemoApp.GetBaseUrl() + MemoApp.URL.TopicsCreateSubmit,
                data: {
                    topic: topicName,
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
     * @param {string} topic
     * @param {jQuery} $allPosts
     * @param {number} lastPostId
     */
    MemoApp.WatchNewTopics = function (topic, $allPosts, lastPostId) {
        $allPosts.scrollTop($allPosts[0].scrollHeight);
        socket = MemoApp.GetSocket(MemoApp.GetBaseUrl() + MemoApp.URL.TopicsSocket + "?topic=" + topic + "&lastPostId=" + lastPostId, function () {
            console.log("socket closed...");
        });
        /**
         * @param {MessageEvent} msg
         */
        socket.onmessage = function (msg) {
            var txHash = msg.data.replace(/['"]+/g, '');
            $.ajax({
                url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoPostAjax + "/" + txHash,
                success: function(html) {
                    $allPosts.append(html);
                    $allPosts.scrollTop($allPosts[0].scrollHeight);
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
    MemoApp.Form.NewTopicMessage = function ($broadcasting, $form) {
        var $topicName = $form.find("[name=topic]");
        var $message = $form.find("[name=message]");

        var $msgByteCount = $form.find(".message-byte-count");

        $message.on("input", function () {
            setMsgByteCount();
        });

        function getByteSize() {
            return MemoApp.utf8ByteLength($topicName.val()) + MemoApp.utf8ByteLength($message.val());
        }

        function setMsgByteCount() {
            var cnt = maxNewTopicBytes - getByteSize();
            $msgByteCount.html("[" + cnt + "]");
            if (cnt < 0) {
                $msgByteCount.addClass("red");
            } else {
                $msgByteCount.removeClass("red");
            }
        }

        var $passwordArea = $form.find("#password-area");
        var $passwordClear = $form.find("#password-clear");

        setMsgByteCount();
        MemoApp.CheckLoadPassword($form);
        if (localStorage.WalletPassword) {
            $passwordArea.hide();
            $passwordClear.addClass("show");
        }
        $passwordClear.find("a").click(function(e) {
            e.preventDefault();
            $passwordArea.show();
            $passwordClear.removeClass("show");
            delete(localStorage.WalletPassword);
            $passwordArea.find("[name=password]").val("");
            $passwordArea.find("[name=save-password]").prop('checked', false);
        });
        $passwordArea.find("[name=save-password]").change(function() {
            if (this.checked) {
                $passwordArea.hide();
                $passwordClear.addClass("show");
                MemoApp.CheckSavePassword($form);
            }
        });
        var submitting = false;
        $form.submit(function (e) {
            e.preventDefault();
            if (submitting) {
                return
            }

            var topicName = $topicName.val();
            var message = $message.val();
            if (maxNewTopicBytes - getByteSize() < 0) {
                alert("Maximum size is " + maxNewTopicBytes + " bytes. Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }

            if (topicName.length === 0) {
                alert("Must enter a topic name.");
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
                url: MemoApp.GetBaseUrl() + MemoApp.URL.TopicsCreateSubmit,
                data: {
                    topic: topicName,
                    message: message,
                    password: password
                },
                success: function (txHash) {
                    if (!txHash || txHash.length === 0) {
                        submitting = false;
                        alert("Server error. Please try refreshing the page.");
                        return
                    }
                    $broadcasting.addClass("show");
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: txHash
                        },
                        success: function () {
                            submitting = false;
                            $broadcasting.removeClass("show");
                            $message.val("");
                            setMsgByteCount();
                        },
                        error: function () {
                            submitting = false;
                            alert("Error waiting for transaction to broadcast.");
                            $broadcasting.removeClass("show");
                            $message.val("");
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
