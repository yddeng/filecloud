let winWidth,winHeight;
let mouse;
let pathDom;
let box;

window.onload = function() {
    if (typeof(Worker) !== "undefined") {
        console.log("浏览器支持HTML5");
    } else {
        alert("浏览器不支持HTML5");
    }

    box = document.getElementById('content');
    box.addEventListener("dragover",function (e) { // 拖来拖去
        e.preventDefault();
    },false);
    box.addEventListener("dragleave",function(e) { // 拖离
        e.preventDefault();
    },false);
    box.addEventListener("drop",dropEvent,false); // 扔

    //获取可视区宽度
    winWidth = function(){ return document.documentElement.clientWidth || document.body.clientWidth;};
    //获取可视区高度
    winHeight = function (){ return document.documentElement.clientHeight || document.body.clientHeight;};

    mouse = document.getElementById('mouse-right');
    pathDom = document.getElementById('path');
    pathDom.value = root;
    document.addEventListener('click', function() {
        mouse.style.display = 'none';
    });
    //右键菜单
    document.oncontextmenu = function(event) {
        var event = event || window.event;
        mouse.style.display = 'block';
        var l, t;
        l = event.clientX;
        t = event.clientY;
        if( l >= (winWidth() - mouse.offsetWidth) ) {
            l  = winWidth() - mouse.offsetWidth;
        }
        if(t > winHeight() - mouse.offsetHeight  ) {
            t = winHeight() - mouse.offsetHeight;
        }
        mouse.style.left = l + 'px';
        mouse.style.top = t + 'px';
        return false;
    };
    refresh()
};

function refresh() {
    fileGet(pathDom.value)
}

function newFolder() {
    let tmp = `<div class="item" >
                <div class="item-check"><input type="checkbox" name="checkbox" class="checkbox"" ></div>
                <div class="item-name"><input id="folder-name" type="text" class="form-control" style="max-width: 400px" autofocus="autofocus" onblur="mkdir()" onkeydown="if(event.keyCode==13){mkdir()}"/></div>
                <div class="item-size">-</div>
                <div class="item-date">-</div>
            </div>`;
    let list = document.getElementById('item-list');
    let str = tmp + list.innerHTML;
    list.innerHTML = str;
}

function mkdir() {
    let input = document.getElementById('folder-name');
    if (input.value === ""){
        showTips("文件名不能为空",2000);
        refresh();
        return;
    }

    let reqUrl = httpAddr+"mkdir?path="+pathDom.value+"/"+input.value;
    util.httpGet(reqUrl,function (res) {
        refresh();
        if (res.ok) {
            showTips("成功",2000);
        }else {
            showTips(res.message,2000)
        }
    },function (e) {
        console.log(e);
        refresh();
        showTips("网络错误！",2000)
    })
}

function fileDelete() {
    let chk_list = document.getElementsByName("checkbox");
    for(let i=0;i<chk_list.length;i++){
        if (chk_list[i].checked) {
            let reqUrl = util.format("{0}delete?path={1}&filename={2}",httpAddr,pathDom.value,chk_list[i].value);
            console.log(reqUrl);
            util.httpGet(reqUrl, function (res) {
                if (res.ok) {
                    showTips("成功", 2000);
                    refresh();
                } else {
                    showTips(res.message, 2000)
                }
            }, function (e) {
                console.log(e);
                showTips("网络错误！", 2000)
            })
        }
    }
}

function fileDownload() {
    let chk_list = document.getElementsByName("checkbox");
    for(let i= 0;i<chk_list.length;i++){
        if (chk_list[i].checked) {
            let reqUrl = util.format("{0}download?path={1}&filename={2}",httpAddr,pathDom.value,chk_list[i].value);
            let tt = chk_list[i].getAttribute("data");
            if (tt === "file") {
                // 创建隐藏的可下载链接
                let eleLink = document.createElement('a');
                eleLink.style.display = 'none';
                eleLink.href = reqUrl;
                // 触发点击
                document.body.appendChild(eleLink);
                eleLink.click();
                // 然后移除
                document.body.removeChild(eleLink);
            }
        }
    }
}

function select(id) {
    if (id !== "") {
        checkAllState(false);
        let name = "checkbox-" + id;
        document.getElementById(name).checked = true;
    }
}

function makeNav(path) {
    let tmp = `<li class="breadcrumb-item"><a href="#" onclick="fileGet('{1}')">{0}</a></li>`;
    let tmpActive = `<li class="breadcrumb-item active" aria-current="page">{0}</li>`;
    let dom = document.getElementById('path-nav');
    dom.innerHTML = "";
    let array = path.split("/");
    let str = "";
    let nowPath = "";
    for (let i = 0;i < array.length;i++){
        if (i === 0) {
            nowPath +=  array[i]
        }else {
            nowPath += "/" + array[i]
        }
        if (i === array.length - 1){
            str += util.format(tmpActive,array[i])
        }else {
            str += util.format(tmp,array[i],nowPath)
        }
    }
    dom.innerHTML = str;
}

function fileGet(path) {
    let tmp = `<div class="item">
                <div class="item-check"><input type="checkbox" id="checkbox-{0}" value="{0}" data="{4}" name="checkbox" class="checkbox" onclick="checkState(this.checked)" ></div>
                <div class="item-name" onmousedown="select('{0}')">{1}</div>
                <div class="item-size" onmousedown="select('{0}')">{2}</div>
                <div class="item-date" onmousedown="select('{0}')">{3}</div>
            </div>`;
    pathDom.value = path;
    let list = document.getElementById('item-list');
    list.innerHTML = "";
    let reqUrl = httpAddr+"list?path="+path;
    makeNav(path);

    let disk_p = document.getElementById('disk-progress');
    let disk_info = document.getElementById('disk-info');

    console.log(reqUrl,path);
    util.httpGet(reqUrl,function (res) {
        if (res.ok) {
            let str = "";
            for (let key in res.items){
                let data = res.items[key];
                if (data.is_dir){
                    let folder = util.format(`<a href="#" onclick="fileGet('{1}')" style="text-decoration: none"><i class="fa fa-folder-o"></i><span>&nbsp;&nbsp;{0}</span></a>`,data.filename,pathDom.value+"/"+data.filename);
                    str += util.format(tmp,data.filename,folder,"-","-","dir")
                }else {
                    let file = util.format(`<i class="fa fa-file-o"></i><span>&nbsp;&nbsp;{0}</span>`,data.filename);
                    str += util.format(tmp,data.filename,file,util.b2string(data.size,1000),data.date,"file")
                }
            }
            list.innerHTML = str;

            disk_info.innerHTML = util.format("{0}/{1}&nbsp;",util.b2string(res.disk_used,1000),util.b2string(res.disk_total,1000))
            disk_p.innerHTML = util.format("{0}%",res.disk_used_p.toFixed(1));
            disk_p.style.width = util.format("{0}%",res.disk_used_p.toFixed(1));
        } else {
            showTips("请求错误", 1000)
        }
    },function (e) {
        console.log(11,e);
        showTips("网络错误！",1000)
    });
}

/************* 文件上传 ***************/

function dropEvent(e) {
    e.preventDefault(); //取消默认浏览器拖拽效果
    let fileList = e.dataTransfer.files; //获取文件对象
    //检测是否是拖拽文件到页面的操作
    if (fileList.length === 0) {
        return false;
    }
    console.log(fileList.length);

    for (let i = 0;i < fileList.length;i++) {
        fileRead(pathDom.value, fileList[i]);
    }
}

function fileRead(path,file) {
    let fileReader = new FileReader();
    fileReader.readAsDataURL(file.slice(0, 4));
    fileReader.onload = function (ev) {
        updateFile(path,file)
    };
    fileReader.onerror = function (ev) {
        console.log(file.name);
        showTips("文件夹"+ file.name,1000)
    };
}

function inputFile() {
    let fileList = document.getElementById('upfile').files;
    if (fileList.length === 0) {
        return
    }
    for (let i = 0;i < fileList.length;i++) {
        updateFile(pathDom.value, fileList[i]);
    }
    document.getElementById('upfile').value = '';
}

let updateInfos = new Map();

function updateInfoClose() {
    updateInfos.clear();
    document.getElementById('update-info').style.display = "none";
    document.getElementById('update-list').innerHTML = "";
}

function updateFile(path,file) {
    let tmp =`<div class="progress">
            <div class="progress-info">
            <div class="progress-name">{0}</div>
            <div class="progress-size">{1}</div>
            <div class="progress-path">{2}</div>
            <div id="{3}" class="progress-width">0%</div>
            </div>
            <div id="{4}" class="progress-bar progress-bar-striped" style="width: 0%;"  aria-valuemin="0" aria-valuemax="100"></div>
        </div>`;

    let chunkSize = 1024 * 1024 * 2, // 以每片2MB大小来逐次读取
        totalSize = file.size,
        filename = file.name,
        fileAbs = path+"/"+filename,
        total = Math.ceil(file.size / chunkSize),
        existBlob = new Map();

    console.log(filename,totalSize,total,fileAbs);
    function addUpdateInfo() {
        document.getElementById('update-info').style.display = "block";
        if (!updateInfos.has(fileAbs)) {
            let list = document.getElementById('update-list');
            let pSize = util.b2string(totalSize,1000);
            let index = path.lastIndexOf("\/");
            let pPath = path.substring(index+1,path.length);
            let pWidth = util.format("{0}-width",fileAbs);
            let pBar = util.format("{0}-bar",fileAbs);
            list.innerHTML += util.format(tmp, filename,pSize,pPath,pWidth,pBar);
            updateInfos.set(fileAbs, fileAbs)
        }
    }

    function updateInfo(sendCnt,isSend) {
        if (updateInfos.has(fileAbs)) {
            let p = (sendCnt*100)/total;
            let ps = util.format("{0}%",p.toFixed(0));
            let pWidth = util.format("{0}-width",fileAbs);
            let pBar = util.format("{0}-bar",fileAbs);
            if (p === 100) {
                document.getElementById(pBar).style.width = "0%";
                if (isSend) {
                    document.getElementById(pWidth).innerHTML = `<i class="fa fa-check-circle" style="color: #00CC00"></i>`
                }else {
                    document.getElementById(pWidth).innerHTML = `<i class="fa fa-check-circle" style="color: #00CC00">秒传</i>`
                }
            }else {
                document.getElementById(pWidth).innerHTML = ps;
                document.getElementById(pBar).style.width = ps;
            }
        }

    }

    function md5File(file,callback) {
        let spark = new SparkMD5(),
            fileReader = new FileReader();
        fileReader.readAsBinaryString(file);
        fileReader.onload = function (ev) {
            spark.appendBinary(ev.target.result);
            callback(spark.end());
        };
    }

    function checkFile(md5,callback) {
        let reqUrl = httpAddr +"check";
        let cmd = {path:path,filename:filename,total:total,md5:md5,size:totalSize};
        util.httpPost(reqUrl,JSON.stringify(cmd), function (res) {
            if (res.ok){
                if (res.need){
                    callback(true,res.upload)
                }else {
                    callback(false)
                }
            }else {
                showTips(res.message, 2000);
            }
        })
    }

    function updateFileEnd() {
        refresh();
    }


    md5File(file,function (md5) {
        checkFile(md5,function (need,exist) {
            addUpdateInfo();
            if (need){
                if (exist) {
                    console.log(exist);
                    for (let key in exist){
                        existBlob.set(key, key)
                    }
                }

                let reqUrl = httpAddr +"upload";
                let current = 0;
                while (current < total) {
                    current++;
                    if (!existBlob.has(current.toString())){
                        let start = chunkSize * (current-1);
                        let end = chunkSize * current;
                        end = end > totalSize ? totalSize : end;
                        let blob = file.slice(start, end);

                        let fd = new FormData();
                        fd.append('path',path);
                        fd.append('file', blob);
                        fd.append('filename', filename);
                        fd.append('current', current.toString());
                        fd.append('md5', md5);

                        util.httpFormData(reqUrl,fd,function (res) {
                            if (res.ok){
                                existBlob.set(fd.get("current"),"0");
                                updateInfo(existBlob.size,true);
                                if (existBlob.size === total){
                                    updateFileEnd()
                                }
                            }else {
                                showTips(res.message,2000)
                            }
                        },function (e) {
                            showTips("网络错误！",1000)
                        });
                    }
                }

            }else {
                updateInfo(total,false);
                updateFileEnd();
            }
        })
    });

}
