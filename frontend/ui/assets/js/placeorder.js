$(document).ready(function() {
    // let items = document.getElementsByClassName("items");
    // let strids = [];
    // for(let x=0; x < items.length; x++){
    //     let itemid = ($(this).children('h4')[0])[x];
    //     strids.push(itemid);
    // }
    // let combined = strids.join(", ");
    let useridd = localStorage.getItem('userid');
    let all = $('.items h4').map(function(){
        return $(this).html();
    }).get().join(',');

    let tot = $("#total").text();
    let dataaa = {"userid": useridd, "prodids": all, "total": tot};
    // let dataaa = {"userid": 2, "prodids": all, "total": tot};

    console.log("order user" + dataaa.userid);
    console.log("order prod" + dataaa.prodids);
    console.log("order total" + dataaa.total);
    $("#checkout").click(function(){
        $.ajax({
            type:"POST",
            url:"http://localhost:8080/cart/buy",
            data: dataaa,
            contentType: 'application/x-www-form-urlencoded',
            success: function(res) {
                console.log(res);
                console.log("Added to order");
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(url, status, err.toString());
            }.bind(this)
        });
    });
});
