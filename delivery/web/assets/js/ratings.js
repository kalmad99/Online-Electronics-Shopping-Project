$(document).ready(function(){
    // Check Radio-box
    $(".ratingfix input:radio").attr("checked", false);

    $('.ratingfix input').click(function () {
        $(".ratingfix span").removeClass('checked');
        $(this).parent().addClass('checked');
    });

    $('input:radio').change(
        function(){
            var userRating = this.value;
            // alert(userRating);
            var rating = document.getElementById('rate');
            rating.value = userRating;
        });
});