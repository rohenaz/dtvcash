(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.ViewBlock = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var heightString = $ele.find("[name=height]").val();
            var height = parseInt(heightString);
            if (height === 0) {
                alert("Must enter a numeric height.");
                return;
            }
            window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Block + "/" + height;
        });
    };
})();
