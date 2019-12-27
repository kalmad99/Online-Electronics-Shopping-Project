$(function() {
    if (localStorage.getItem('select')) {
        $("#select option").eq(localStorage.getItem('select')).prop('selected', true);
    }
    localStorage.removeItem('select');

    $("#select").on('change', function() {
        localStorage.setItem('select', $('option:selected', this).index());
    });
});
