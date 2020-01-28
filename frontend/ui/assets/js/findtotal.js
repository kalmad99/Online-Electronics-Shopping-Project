$(document).ready(function(){
    let prices = document.getElementsByClassName("price");
    let total = 0;
    for(let x=0; x < prices.length; x++)
    {
        let intpri = parseFloat((prices[x]).innerHTML);
        total += intpri;
    }
    $("#total").html(total);
});
// console.log(price.innerHTML);
// console.log(total);