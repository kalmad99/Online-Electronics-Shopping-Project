$(document).ready(function() {
    let proid = document.getElementById("id");
    let currentdate = new Date();
    let price = document.getElementById("price");
    let useridd = localStorage.getItem("userid");

    console.log("uId",useridd);
    //ok
    let datetime = currentdate.getFullYear() + "-" +
        (currentdate.getMonth()+1)  + "-" +
        currentdate.getDate() + " "
        + currentdate.getHours() + ":"
        + currentdate.getMinutes() + ":"
        + currentdate.getSeconds();

    // let dataaa = {"prodid": proid.innerHTML, "userid":2, "addedtime":datetime, "price": price.innerHTML};
    let dataaa = {"prodid": proid.innerHTML, "userid":useridd, "addedtime":datetime, "price": price.innerHTML};

    $("#addtocart").click(function(){
        $.ajax({
            type:"POST",
            url:" http://localhost:8080/addtocart",
            data: dataaa,
            contentType: 'application/x-www-form-urlencoded',
            success: function(res) {
                console.log(res);
                console.log("Added");
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(url, status, err.toString());
            }.bind(this)
        });
    });
});
