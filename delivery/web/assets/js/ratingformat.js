$(document).ready(function(){
    let ratevalue = document.getElementById('ratevalue').innerText;
    console.log(ratevalue);
    $("input[name=rating][value=" + ratevalue*2 + "]").attr('checked', 'checked');
});
