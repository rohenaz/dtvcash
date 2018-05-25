(function () {
    var maxOptionBytes = 184;
    var maxQuestionBytes = 209;

    var $form;
    var $question;
    var $options;
    var $questionByteCount;
    var optionCounter = 1;
    var submitting = false;

    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.NewPoll = function ($ele) {
        $form = $ele;
        $options = $form.find("#options");
        $question = $form.find("[name=question]");
        $questionByteCount = $form.find(".question-byte-count");

        $form.find("#add-option").click(function () {
            addOption();
            bindRemoveOption();
        });

        $question.on("input", function () {
            setQuestionByteCount();
        });

        addOption();
        bindRemoveOption();
        initOptionByteCounter(1);
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

        var $optionInputs = $options.find("[name=option]");
        if ($optionInputs.length < 2) {
            alert("Error, not enough options.");
            return;
        }
        var options = [];
        for (var i = 0; i < $optionInputs.length; i++) {
            var option = $optionInputs.eq(i).val();
            if (maxOptionBytes - MemoApp.utf8ByteLength(option) < 0) {
                alert("Maximum option size is " + maxOptionBytes + " bytes." +
                    " Note that some characters are more than 1 byte." +
                    " Emojis are usually 4 bytes, for example.");
                return;
            }
            options.push(option);
        }

        postPoll(pollType, question, options, password);
    }

    /**
     * @param {string} pollType
     * @param {string} question
     * @param {[string]} options
     * @param {string} password
     */
    function postPoll(pollType, question, options, password) {
        submitting = true;
        $.ajax({
            type: "POST",
            url: MemoApp.GetBaseUrl() + MemoApp.URL.PollCreateSubmit,
            data: {
                pollType: pollType,
                question: question,
                options: options,
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
                    "Error with request (option code " + xhr.status + "):\n" +
                    (xhr.optionText !== "" ? xhr.optionText + "\n" : "") +
                    "If this problem persists, try refreshing the page.";
                alert(errorMessage);
            }
        });
    }

    function addOption() {
        var $firstOption = $options.find("div:eq(0)");
        var newOptionHtml = $firstOption.clone()[0].outerHTML.replace(/option-[0-9]+/g, "option-" + ++optionCounter);
        $options.append(newOptionHtml);
        initOptionByteCounter(optionCounter);
        showHideRemoveButton();
    }

    /**
     * @param {number} id
     */
    function initOptionByteCounter(id) {
        var $option = $options.find("#option-" + id);
        var $optionByteCounter = $option.find(".option-byte-count");
        var $input = $option.find("input");
        $input.on("input", function () {
            setOptionByteCounter();
        });

        function setOptionByteCounter() {
            var cnt = maxOptionBytes - MemoApp.utf8ByteLength($input.val());
            $optionByteCounter.html("[" + cnt + "]");
            if (cnt < 0) {
                $optionByteCounter.addClass("red");
            } else {
                $optionByteCounter.removeClass("red");
            }
        }

        setOptionByteCounter();
    }

    function bindRemoveOption() {
        $options.find(".remove-option").click(function (e) {
            e.preventDefault();
            var numOptions = $options.find("> div").length;
            if (numOptions <= 2) {
                return;
            }
            var $this = $(this);
            var optionId = $this.attr("data-option-id");
            var $option = $("#" + optionId);
            $option.remove();
            showHideRemoveButton();
        });
    }

    function showHideRemoveButton() {
        var $removeButtons = $options.find(".remove-option");
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
