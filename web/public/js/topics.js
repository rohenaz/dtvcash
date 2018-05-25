(function () {

    var maxNewTopicBytes = 204;

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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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
     */
    MemoApp.WatchNewTopics = function (topic, $allPosts) {
        function connect() {
            topic = encodeURIComponent(topic);
            var params = "?topic=" + topic + "&lastPostId=" + _lastPostId + "&lastLikeId=" + _lastLikeId;
            var socket = MemoApp.GetSocket(MemoApp.GetBaseUrl() + MemoApp.URL.TopicsSocket + params);

            socket.onclose = function () {
                setTimeout(function () {
                    connect();
                }, 1000);
            };
            /**
             * @param {MessageEvent} msg
             */
            socket.onmessage = function (msg) {
                var data;
                try {
                    data = JSON.parse(msg.data);
                } catch (e) {
                    return;
                }
                var txHash = data.Hash.replace(/['"]+/g, '');
                var $post = $("#topic-post-" + txHash);
                if (data.Type === 2 && !$post.length) {
                    return;
                }
                $.ajax({
                    url: MemoApp.GetBaseUrl() + MemoApp.URL.TopicsPostAjax + "/" + txHash,
                    success: function (html) {
                        if ($post.length) {
                            $post.replaceWith(html);
                            return;
                        }
                        $allPosts.append(html);
                        $allPosts.scrollTop($allPosts[0].scrollHeight);
                    },
                    error: function (xhr) {
                        alert("error getting post via ajax (status: " + xhr.status + ")");
                    }
                });
            };
        }

        connect();
    };

    /**
     * @param {string} topic
     * @param {jQuery} $allPosts
     */
    MemoApp.LoadMore = function (topic, $allPosts) {
        var submitting = false;
        $allPosts.scroll(function () {
            if (submitting) {
                return;
            }
            var pos = $allPosts.scrollTop();
            if (pos === 0) {
                submitting = true;
                $.ajax({
                    url: MemoApp.GetBaseUrl() + MemoApp.URL.TopicsMorePosts,
                    data: {
                        firstPostId: _firstPostId,
                        topic: topic
                    },
                    success: function (html) {
                        submitting = false;
                        if (html === "") {
                            return;
                        }
                        var firstItem = $allPosts.find(":first");
                        var curOffset = firstItem.offset().top - $allPosts.scrollTop();
                        $allPosts.prepend(html);
                        $allPosts.scrollTop(firstItem.offset().top - curOffset);
                    },
                    error: function (xhr) {
                        submitting = false;
                        alert("error getting posts (status: " + xhr.status + ")");
                    }
                });
            }
        });
    };

    var _firstPostId;
    var _lastPostId;
    var _lastLikeId;

    /**
     * @param {number} firstPostId
     */
    MemoApp.SetFirstPostId = function (firstPostId) {
        if (_firstPostId === undefined || firstPostId < _firstPostId) {
            _firstPostId = firstPostId;
        }
    };

    /**
     * @param {number} lastPostId
     */
    MemoApp.SetLastPostId = function (lastPostId) {
        if (_lastPostId === undefined || lastPostId > _lastPostId) {
            _lastPostId = lastPostId;
        }
    };

    /**
     * @param {number} lastLikeId
     */
    MemoApp.SetLastLikeId = function (lastLikeId) {
        if (_lastLikeId === undefined || lastLikeId > _lastLikeId) {
            _lastLikeId = lastLikeId;
        }
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

        setMsgByteCount();
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

            var password = MemoApp.GetPassword();
            if (!password.length) {
                console.log("Password not set. Please try logging in again.");
                return;
            }

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
                    $broadcasting.removeClass("hidden");
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: txHash
                        },
                        success: function () {
                            submitting = false;
                            $broadcasting.addClass("hidden");
                            $message.val("");
                            setMsgByteCount();
                        },
                        error: function () {
                            submitting = false;
                            alert("Error waiting for transaction to broadcast.");
                            $broadcasting.addClass("hidden");
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

    /**
     * @param {jQuery} $like
     * @param {string} txHash
     */
    MemoApp.Form.NewTopicLike = function ($like, txHash) {
        $like.find("#like-link-" + txHash).click(function (e) {
            e.preventDefault();
            $("#like-info-" + txHash).hide();
            $("#like-form-" + txHash).css({"display": "inline-block"});
        });
        $like.find("#like-cancel-" + txHash).click(function (e) {
            e.preventDefault();
            $("#like-info-" + txHash).show();
            $("#like-form-" + txHash).css({"display": "none"});
        });
        var $form = $like.find("form");

        var $broadcasting = $like.find(".broadcasting");

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
                    $broadcasting.removeClass("hidden");
                    $form.hide();
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: txHash
                        },
                        success: function () {
                            submitting = false;
                            $broadcasting.addClass("hidden");
                        },
                        error: function () {
                            submitting = false;
                            alert("Error waiting for transaction to broadcast.");
                            $broadcasting.addClass("hidden");
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
