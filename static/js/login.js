function login(){
    let username = document.getElementById("loginusername")
    let password = document.getElementById("loginpassword")
    var url = '/v1/account/login';
    console.log(password)
    let data = {
        username: username.value,
        password: password.value
    }
    fetch(url, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data), // data can be `string` or {object}!
        headers: new Headers({
          'Content-Type': 'application/json'
        })
    }).then(res => {
        if ( res.status == 200 ) {
            window.location.replace("/");
        }else{
            throw new Error("no login")
        }
    }).catch(()=>{
        let toast = document.getElementById("toast")
        toast.innerHTML = "登入失敗"
    })
}

function createAccount(){
    let username = document.getElementById("createusername")
    let password = document.getElementById("createpassword")
    var url = '/v1/account/registry';
    console.log(password)
    let data = {
        username: username.value,
        password: password.value
    }
    fetch(url, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data), // data can be `string` or {object}!
        headers: new Headers({
          'Content-Type': 'application/json'
        })
    }).then(res => {
        if ( res.status == 200 ) {  
            let toast = document.getElementById("toast")
            toast.innerHTML = "建置成功"
        }else{
            throw new Error("no registry")
        }
    }).catch(()=>{
        let toast = document.getElementById("toast")
        toast.innerHTML = "建置失敗"
    })
}