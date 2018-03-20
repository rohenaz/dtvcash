(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.ViewBlock = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var heightString = $ele.find("[name=height]").val();
            var height = parseInt(heightString);
            if (height < 0 || heightString != height) {
                alert("Must enter a numeric height.");
                return;
            }
            window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Block + "/" + height;
        });
    };

    /**
     * @param {jQuery} $nextBtn
     * @param {jQuery} $prevBtn
     */
    MemoApp.Form.BindBlockArrows = function($nextBtn, $prevBtn) {
        $(document).keydown(function(e) {
            switch(e.which) {
                case 37:
                    $prevBtn[0].click();
                    break;
                case 39:
                    $nextBtn[0].click();
                    break;
                default: return;
            }
            e.preventDefault();
        });
    };
})();
