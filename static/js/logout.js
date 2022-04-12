function logout(){
    var url = '/v1/account/logout';
    fetch(url, {
        method: 'GET', // or 'PUT'
    }).then(res => {
        if ( res.status == 200 ) {
            window.location.replace("/");
        }else{
            throw new Error("no login")
        }
    })
}