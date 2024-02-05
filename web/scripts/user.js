let email = document.querySelector("#emaildiv");
let apikey = document.querySelector("#apikeydiv");
let api = "";

function GetUserInfo() {
    $.ajax({
        url: "http://127.0.0.1:8080/user/",
        method: "GET",
        cache: false,
        crossDomain: true,
        xhrFields: { withCredentials: true },
        success : function (data, textStatus, xhr) {

           loadPage(data);

        },
        error: function(data, textStatus, xhr) { 
            console.log(data.status, textStatus, xhr.status);
        } 
        });
}

function loadPage(data) {
    let tunneldiv = document.querySelector("#tunneldiv");
    console.log(data);
    email.textContent = data.login;
    api = data.apikey;
    apikey.textContent = data.apikey;
    // console.log(data.Tunnels.length)

    if (data.Tunnels.length==0) {
        tunneldiv.innerHTML="<h5 class=\"font-monospace mt-5\">There are no configured tunnels</h5>";
    }  else {
        tunneldiv.toggle(); 
    }
}

addtunnel.onsubmit = async (e) => {
    e.preventDefault();
    var form = document.querySelector("#addtunnel");
    let ipv4client = form.querySelector('#ipv4client').value;
    let tunnelname = form.querySelector('#tunnelname').value;
    data = {
        tunnelname : tunnelname,
        ipv4remote : ipv4client,
    };
    $.ajax({
    url: "http://127.0.0.1:8080/tunnel/"+api+"/",
    method: "POST",
    cache: false,
    data: JSON.stringify(data),
    crossDomain: true,
    xhrFields: { withCredentials: true },
    headers: {
        'Content-Type': 'application/json',
    },
    success : function (data, textStatus, xhr) {
        console.log(data, textStatus, xhr);
        // if (xhr.status==200) {
        //     // xhr.getResponseHeader('Set-Cookie');
        //     result.textContent = "Success";
        //     var count = 1;
        //     setInterval(function(){
        //         count--;
        //         if (count == 0) {
        //             window.location = '/user.html'; 
        //         }
        //     },1000);
        // }   
    },
    error: function(data, textStatus, xhr) {
        result.textContent = "email or password incorrect"
        console.log(data.status, textStatus, xhr.status);
    } 
    });
}
