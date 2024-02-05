useraddform.onsubmit = async (e) => {
    e.preventDefault();
    var form = document.querySelector("#useraddform");
    let email = form.querySelector('#email').value;
    let password = form.querySelector('#password').value;
    let rpassword = form.querySelector('#repeatpassword').value;
    let result = document.querySelector("#result");
    let submit = document.querySelector("#createbtn");
    if (password != rpassword) {
        result.textContent="Passwords don't match"
        return;
    } else {
        result.textContent="";
        submit.disabled = true
    }
    data = {
        login : email,
        password : password,
    };
    $.ajax({
    url: "http://127.0.0.1:8080/signup",
    method: "POST",
    data: JSON.stringify(data),
    headers: {
        'Content-Type': 'application/json',
    },
    success : function (data, textStatus, xhr) {
        console.log(data.status, textStatus, xhr.status);
        result.innerHTML="User is created.<br>You will be redirected to login page in few seconds <div id=countDown></div>";
        var count = 5;
        setInterval(function(){
            count--;
            document.getElementById('countDown').innerHTML = count;
            if (count == 0) {
                window.location = '/index.html'; 
            }
        },1000);
    },
    error: function(data, textStatus, xhr) { 
        console.log(data.status, textStatus, xhr.status);
        submit.disabled = false;
    } 
    });
}
