(function() {

    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.NewPoll = function ($form) {
        var $addResponse = $form.find("#add-response");
        var $responses = $form.find("#responses");
        var responseCounter = 1;

        $addResponse.click(function() {
            addResponse();
            bindRemoveResponse();
        });

        addResponse();
        bindRemoveResponse();

        function addResponse() {
            var $firstResponse = $responses.find("div:eq(0)");
            var newResponseHtml = $firstResponse.clone()[0].outerHTML.replace(/response-[0-9]+/g, "response-" + ++responseCounter);
            $responses.append(newResponseHtml);
        }

        function bindRemoveResponse() {
            $responses.find(".remove-response").click(function(e) {
                e.preventDefault();
                var numResponses = $responses.find("> div").length;
                if (numResponses <= 2) {
                    return;
                }
                var $this = $(this);
                var responseId = $this.attr("data-response-id");
                var $response = $("#" + responseId);
                $response.remove();
            });
        }
    };
})();
