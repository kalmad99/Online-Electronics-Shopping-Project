$(document).ready(function() {
    let usrid = localStorage.getItem("userid");
    if (userid == ""){
        $('a[href*="/login"]').attr('href', '/login');
    }else {
        $('a[href*="/login"]').attr('href', '/userprof?id='+usrid);
    }
});
