
function regbtn(){
    window.location = '/register.html'; 
}


loginform.onsubmit = async (e) => {
    e.preventDefault();
    var form = document.querySelector("#loginform");
    let email = form.querySelector('#email').value;
    let password = form.querySelector('#password').value;
    let result = document.querySelector("#result");
    let login = document.querySelector("#loginbtn");
    data = {
        login : email,
        password : password,
    };
    $.ajax({
    url: "http://127.0.0.1:8080/login",
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
        if (xhr.status==200) {
            // xhr.getResponseHeader('Set-Cookie');
            result.textContent = "Success";
            var count = 1;
            setInterval(function(){
                count--;
                if (count == 0) {
                    window.location = '/user.html'; 
                }
            },1000);
        }   
    },
    error: function(data, textStatus, xhr) {
        result.textContent = "email or password incorrect"
        console.log(data.status, textStatus, xhr.status);
    } 
    });
}

