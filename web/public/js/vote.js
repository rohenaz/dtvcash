(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.NewVote = function ($ele) {
        var $voteShowForm = $ele.find(".vote-show-form");
        var $results = $ele.find(".results");
        var $form = $ele.find("form");
        var $cancelButton = $ele.find(".vote-cancel");
        $voteShowForm.click(function (e) {
            e.preventDefault();
            $form.removeClass("hidden");
            $results.addClass("hidden");
        });
        $cancelButton.click(function (e) {
            e.preventDefault();
            $form.addClass("hidden");
            $results.removeClass("hidden");
        });
    };
})();
