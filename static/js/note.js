function actionTd(id){
    return (
        "<td>"
        + "<button class=\"btn\" data-id="+id+" onclick=\"deleteNote(this)\">刪除</button>"
        + "<button class=\"btn\" data-id="+id+" onclick=\"showEditNoteModel(this,'PUT')\">修改</button>"
        + "</td>"
    )

}

function getAllNote(){
    url = "/v1/note"
    fetch(url, {
        method: 'GET', // or 'PUT'
    }).then(res => {
        return res.json()
    }).then(jObj => {
        let noteTableBody = document.getElementById("noteTableBody")
        let htmlStr = ""
        for(let obj of jObj){
            let time = new Date(obj.Time)
            htmlStr += "<tr>"
            htmlStr += "<td>" + time.toLocaleDateString('zh-Hans-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }) + "</td>"
            htmlStr += "<td>" + time.toLocaleTimeString() + "</td>"
            htmlStr += "<td>" + obj.Message + "</td>"
            htmlStr += actionTd(obj.Id)
            htmlStr += "</tr>"
        }
        noteTableBody.innerHTML = htmlStr
    }).catch((e)=>{
        console.log(e)
    })
}

function deleteNote(ref){
    let id = ref.getAttribute("data-id")
    console.log(id)
    url = "/v1/note/fileId/"+id
    fetch(url, {
        method: 'DELETE', // or 'PUT'
    }).then(res => {
        if(res.status == 200){
            getAllNote()
            return 
        }else{
            throw new Error("not delete")
        }
    }).catch((e)=>{
        console.log(e)
    })
}

function showEditNoteModel(ref,method){
    let id = ""
    if(method == 'PUT'){
        id = ref.getAttribute("data-id")
        url = "/v1/note/fileId/"+id
        fetch(url, {
            method: 'GET', // or 'PUT'
        }).then(res => {
            if(res.status == 200){
                return res.json()
            }else{
                throw new Error("not GET")
            }
        }).then(jobj => {
            console.log(jobj)
            let time = new Date(jobj.Time)
            time.setMinutes(time.getMinutes() - time.getTimezoneOffset());
            let dateTime = time.toISOString().slice(0,16)
            let message = jobj.Message
            let idDom = document.getElementById("modifyId")
            let modifyDateTimeDom = document.getElementById("modifyDateTime")
            let messageDom = document.getElementById("modifyMessage")
            document.getElementById("btnModify").style.display = "block"
            document.getElementById("btnCreate").style.display = "none"
            idDom.value = id
            modifyDateTimeDom.value = dateTime
            messageDom.value = message
        }).catch((e)=>{
            console.log(e)
        })
    }else if(method == "POST"){
        document.getElementById("modifyId").value = ""
        document.getElementById("modifyDateTime").value = ""
        document.getElementById("modifyMessage").value = ""
        document.getElementById("btnModify").style.display = "none"
        document.getElementById("btnCreate").style.display = "block"
    }
    
}

function modifyNote(method){
    let idDom = document.getElementById("modifyId")
    let dateTimeDom = document.getElementById("modifyDateTime")
    let messageDom = document.getElementById("modifyMessage")
    let id = idDom.value
    let dateTime = new Date(dateTimeDom.value)
    console.log(dateTime)
    let url = "/v1/note"
    if(method == "PUT"){
        url += "/fileId/"+id
    }else{
        
    }
    data = {
        Time: dateTime.toISOString(),
        Message: messageDom.value
    }
    fetch(url, {
        method: method, // or 'PUT'
        body: JSON.stringify(data), // data can be `string` or {object}!
        headers: new Headers({
          'Content-Type': 'application/json'
        })
    }).then(res => {
        if(res.status == 200){
            getAllNote()
            return 
        }else{
            throw new Error("not delete")
        }
    }).catch((e)=>{
        console.log(e)
    })
}