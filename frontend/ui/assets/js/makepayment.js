// $(document).ready(function() {
//     let useridd = localStorage.getItem('userid');
//     let tot = document.getElementById('tot');
//     let dataaa = {"userid": useridd, "total": tot.innerHTML};
//     // let dataaa = {"userid": 2, "prodids": all, "total": tot};
//
//     console.log("bank user" + dataaa.userid);
//     console.log("bank total" + dataaa.total);
//     $("#pay").click(function(){
//         $.ajax({
//             type:"GET",
//             url:"http://localhost:8080/pay",
//             data: dataaa,
//             contentType: 'application/x-www-form-urlencoded',
//             success: function(res) {
//                 console.log(res);
//                 console.log("Added to pay");
//             }.bind(this),
//             error: function(xhr, status, err) {
//                 console.error(url, status, err.toString());
//             }.bind(this)
//         });
//     });
// });
$(document).ready(function(){
    let oid = document.getElementById('oid');
    let userid = localStorage.getItem("userid");
    $("#tryagain").click(function(){
        $.ajax({
            type:"GET",
            url:" http://localhost:8080/order/delete?id=" + oid.innerHTML,
            contentType: 'application/x-www-form-urlencoded',
            success: function(res) {
                console.log(res);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(url, status, err.toString());
            }.bind(this)
        });
        window.location = "/getusercart?id=" + userid;
    });
});