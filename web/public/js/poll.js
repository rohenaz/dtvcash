(function () {
    var maxResponseBytes = 184;
    var maxQuestionBytes = 213;

    var $form;
    var $question;
    var $responses;
    var $questionByteCount;
    var responseCounter = 1;
    var submitting = false;

    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.NewPoll = function ($ele) {
        $form = $ele;
        $responses = $form.find("#responses");
        $question = $form.find("[name=question]");
        $questionByteCount = $form.find(".question-byte-count");

        $form.find("#add-response").click(function () {
            addResponse();
            bindRemoveResponse();
        });

        $question.on("input", function () {
            setQuestionByteCount();
        });

        addResponse();
        bindRemoveResponse();
        initResponseByteCounter(1);
        showHideRemoveButton();
        setQuestionByteCount();

        $form.submit(function (e) {
            e.preventDefault();
            submitForm();
        });
    };

    function submitForm() {
        if (submitting) {
            return
        }

        var pollType = $form.find("[name=poll-type] option:selected").val();
        if (pollType.length === 0) {
            alert("Error getting poll type.");
            return;
        }

        var question = $question.val();
        if (maxQuestionBytes - MemoApp.utf8ByteLength(question) < 0) {
            alert("Maximum question size is " + maxQuestionBytes + " bytes." +
                " Note that some characters are more than 1 byte." +
                " Emojis are usually 4 bytes, for example.");
            return;
        }

        if (question.length === 0) {
            alert("Must enter a question.");
            return;
        }

        var password = MemoApp.GetPassword();
        if (!password.length) {
            console.log("Password not set. Please try logging in again.");
            return;
        }

        var $responseInputs = $responses.find("[name=response]");
        if ($responseInputs.length < 2) {
            alert("Error, not enough responses.");
            return;
        }
        var responses = [];
        for (var i = 0; i < $responseInputs.length; i++) {
            var response = $responseInputs.eq(i).val();
            if (maxResponseBytes - MemoApp.utf8ByteLength(response) < 0) {
                alert("Maximum response size is " + maxResponseBytes + " bytes." +
                    " Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }
            responses.push(response);
        }

        postPoll(pollType, question, responses, password);
    }

    /**
     * @param {string} pollType
     * @param {string} question
     * @param {[string]} responses
     * @param {string} password
     */
    function postPoll(pollType, question, responses, password) {
        submitting = true;
        $.ajax({
            type: "POST",
            url: MemoApp.GetBaseUrl() + MemoApp.URL.PollCreateSubmit,
            data: {
                pollType: pollType,
                question: question,
                responses: responses,
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
    }

    function addResponse() {
        var $firstResponse = $responses.find("div:eq(0)");
        var newResponseHtml = $firstResponse.clone()[0].outerHTML.replace(/response-[0-9]+/g, "response-" + ++responseCounter);
        $responses.append(newResponseHtml);
        initResponseByteCounter(responseCounter);
        showHideRemoveButton();
    }

    /**
     * @param {number} id
     */
    function initResponseByteCounter(id) {
        var $response = $responses.find("#response-" + id);
        var $responseByteCounter = $response.find(".response-byte-count");
        var $input = $response.find("input");
        $input.on("input", function () {
            setResponseByteCounter();
        });

        function setResponseByteCounter() {
            var cnt = maxResponseBytes - MemoApp.utf8ByteLength($input.val());
            $responseByteCounter.html("[" + cnt + "]");
            if (cnt < 0) {
                $responseByteCounter.addClass("red");
            } else {
                $responseByteCounter.removeClass("red");
            }
        }

        setResponseByteCounter();
    }

    function bindRemoveResponse() {
        $responses.find(".remove-response").click(function (e) {
            e.preventDefault();
            var numResponses = $responses.find("> div").length;
            if (numResponses <= 2) {
                return;
            }
            var $this = $(this);
            var responseId = $this.attr("data-response-id");
            var $response = $("#" + responseId);
            $response.remove();
            showHideRemoveButton();
        });
    }

    function showHideRemoveButton() {
        var $removeButtons = $responses.find(".remove-response");
        if ($removeButtons.length > 2) {
            $removeButtons.removeClass("hidden");
        } else {
            $removeButtons.addClass("hidden");
        }
    }

    function setQuestionByteCount() {
        var cnt = maxQuestionBytes - MemoApp.utf8ByteLength($question.val());
        $questionByteCount.html("[" + cnt + "]");
        if (cnt < 0) {
            $questionByteCount.addClass("red");
        } else {
            $questionByteCount.removeClass("red");
        }
    }
})();
