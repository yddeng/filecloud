var util = {};
util.httpGet = function(url,success,error){
    $.ajax({
        url:url,
        type: "get",
        async: true,
        success: success,
        error: error
    });
};

util.httpPost = function(url,data,success,error){
    $.ajax({
        url:url,
        type: "post",
        async: true,
        dataType: "json",
        data:data,
        success: success,
        error: error
    });
};

util.httpFormData = function(url,data,success,error){
    $.ajax({
        url:url,
        type: "post",
        data:data,
        contentType: false,//"multipart/form-data",
        processData: false,
        success: success,
        error: error
    });
};

// 字符串格式化
util.format = function(src){
    if (arguments.length == 0) return null;
    let args = Array.prototype.slice.call(arguments, 1);
    return src.replace(/\{(\d+)\}/g, function(m, i){
        return args[i];
    });
};

// 设置cookie的函数  （名字，值，过期时间（天））
util.setCookie = function (cname, cvalue, exdays) {
    let d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    let expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
};

//获取cookie
//取cookie的函数(名字) 取出来的都是字符串类型 子目录可以用根目录的cookie，根目录取不到子目录的 大小4k左右
util.getCookie = function(cname) {
    let name = cname + "=";
    let ca = document.cookie.split(';');
    for(let i=0; i<ca.length; i++)
    {
        let c = ca[i].trim();
        if (c.indexOf(name)===0) return c.substring(name.length,c.length);
    }
    return "";
};

util.percent = function (v) {
    let n = parseFloat(v);
    return util.format("{0}%", n.toFixed(2))
};

util.str2Int = function (v) {
    return parseInt(v)
};

util.unit = new Array("B", "KB", "MB", "GB");

util.b2string = function(v,rate) {
    let n = v;
    let i = 0;
    for (;n > rate;) {
        n /= rate;
        i++;
        if (i === util.unit.length){
            break
        }
    }
    return util.format("{0}{1}",Math.round(n), util.unit[i])
};

function showTips(msg,t) {
    let tip = document.getElementById('tips');
    tip.style.display = "block";

    let m = document.getElementById("tips-msg");
    m.innerText=msg;
    setTimeout(function(){ tip.style.display = "none"},t)
}

// 全部选中，checkAll状态为true。 有一个选中批量操作按钮可点击，否则不能
function checkState(isChecked) {
    let chk_list = document.getElementsByName("checkbox");
    let checkedNum = 0;
    for(let i=0;i<chk_list.length;i++){
        if (chk_list[i].checked){
            checkedNum++
        }
    }

    if (checkedNum === chk_list.length){
        document.getElementById("checkAll").checked = true;
    }else {
        document.getElementById("checkAll").checked = false;
    }

    let batchBtn = document.getElementsByClassName("batch-btn");
    if (checkedNum === 0){
        for (let i =0;i < batchBtn.length;i++){
            batchBtn[i].disabled = true
        }
    }else {
        for (let i =0;i < batchBtn.length;i++){
            batchBtn[i].disabled = false
        }
    }

}

function checkAllState(isChecked){
    let chk_list = document.getElementsByName("checkbox");
    for(let i=0;i<chk_list.length;i++){
        chk_list[i].checked=isChecked;
    }

    let batchBtn = document.getElementsByClassName("batch-btn");
    if (!isChecked){
        for (let i =0;i < batchBtn.length;i++){
            batchBtn[i].disabled = true
        }
    }else {
        for (let i =0;i < batchBtn.length;i++){
            batchBtn[i].disabled = false
        }
    }
}