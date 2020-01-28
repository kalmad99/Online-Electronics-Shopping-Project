let cookie = document.getElementById("cookiename");
let token = cookie.innerHTML;
let cookieval = $.cookie(token);
// console.log(parseJwt(token));
// let par = parseJwt(csrf.innerHTML);
// if (par){
//     console.log(par)
// }
// function parseJwt(token) {
//     try {
//         // Get Token Header
//         const base64HeaderUrl = token.split('.')[0];
//         const base64Header = base64HeaderUrl.replace('-', '+').replace('_', '/');
//         const headerData = JSON.parse(window.atob(base64Header));
//
//         // Get Token payload and date's
//         const base64Url = token.split('.')[1];
//         const base64 = base64Url.replace('-', '+').replace('_', '/');
//         const dataJWT = JSON.parse(window.atob(base64));
//         dataJWT.header = headerData;
//
// // TODO: add expiration at check ...
//
//
//         return dataJWT;
//     } catch (err) {
//         return false;
//     }
// }
// let csrf = document.getElementById("csrf");
// const jwtDecoded = parseJwt(csrf.innerHTML) ;
// if(jwtDecoded)
// {
//     console.log(jwtDecoded)
// }
// else {
//     console.log("Hello")
// }

// var token = getelemen
let b64DecodeUnicode = str =>
    decodeURIComponent(
        Array.prototype.map.call(atob(str), c =>
            '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
        ).join(''));

let parseJwt = token =>
    JSON.parse(
        b64DecodeUnicode(
            token.split('.')[1].replace('-', '+').replace('_', '/')
        )
    );


// let csrf = document.getElementById("csrf");
let out = document.getElementById("email");
out.value = JSON.stringify(
    parseJwt(cookieval)
);
let email = out.value[0];
//
// console.log(csrf.innerHTML);
console.log(email);


// function parseJwt () {
//     var base64Url = token.split('.')[1];
//     var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
//     var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
//         return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
//     }).join(''));
//
//     return JSON.parse(jsonPayload);
// };