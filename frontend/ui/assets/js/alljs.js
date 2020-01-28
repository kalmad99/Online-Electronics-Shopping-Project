$(document).ready(function() {
    let usrid = localStorage.getItem("userid");
    $('a[href*="#"]').attr('href' , '/getusercart?id=' + usrid);
});