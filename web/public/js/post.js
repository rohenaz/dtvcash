$( document ).ready(function() {

    $("#message").on("keyup", function(e){
        $("#count").html(75 - MemoApp.utf8ByteLength($(e.currentTarget).val()));
    })
})