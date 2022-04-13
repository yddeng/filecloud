<template>
 <div class="body">
  <div id="header">
    <a-row justify="space-between" type="flex">
      <a-col>
        <template v-if="this.selectedNames.length === 0">
          <a-upload 
          ref="upload"
          multiple 
          :directory="uploadDirectory"
          :openFileDialogOnClick="uploadFileOpen"
          :showUploadList="false"
          :beforeUpload="handleBeforeUpload"
          @change="handleChange">
            <a-dropdown>
              <a-button type="primary" shape="round" icon="upload">上传</a-button>
              <a-menu slot="overlay" @click="handleUploadType">
                <a-menu-item key="file">
                  <a-icon type="upload"/> 上传文件
                </a-menu-item>
                 <a-menu-item key="folder">
                  <a-icon type="folder"/> 上传文件夹
                </a-menu-item>
              </a-menu>
            </a-dropdown>
          </a-upload>&nbsp;
          <a-button-group >
            <a-button icon="folder-add" @click="openAddFolder">新建文件夹</a-button>
            <a-button icon="download" @click="()=>{$message.info('功能未实现')}">离线下载</a-button>
            <a-button icon="sync" @click="() => { this.goto(this.navs.length-1) }">刷新</a-button>
          </a-button-group>
        </template>
        <template v-else>
          <a-button-group >
            <a-button icon="share-alt" @click="openShareModal">分享</a-button>
            <a-button icon="download" v-show="this.showDownload()" @click="handleDownload">下载</a-button>
            <a-button icon="delete" @click="openRemoveModal">删除</a-button>
            <a-button icon="form" v-show="this.selectedNames.length === 1" @click="openRename">重命名</a-button>
            <a-button icon="copy" @click="openMvcpModal(false)">复制</a-button>
            <a-button icon="drag" @click="openMvcpModal(true)">移动</a-button>
          </a-button-group>
        </template>
      </a-col>
      <a-col style="width:300px;margin-right:100px">
        <a-row >
          <a-col :span="10">
            {{ resData.diskUsedStr }} / {{ resData.diskTotalStr }}
          </a-col>
          <a-col :span="13">
            <a-progress v-if="resData.diskUsed < resData.diskTotal" 
            :stroke-color="progressColor(resData.diskUsed / resData.diskTotal * 100)"
            :percent="parseFloat((resData.diskUsed / resData.diskTotal * 100).toFixed(1))" />
            <a-progress v-else :percent="100" stroke-color="red" :show-info="false"/>
          </a-col>
        </a-row>
      </a-col>
    </a-row>
  </div>

  <div class="path">
    <template v-if="navs.length === 1">
      全部文件
    </template>
    <template v-else>
      <a href="javascript:;" @click="goback()">返回上一级</a>
      <a-divider type="vertical" />
      <a href="javascript:;" @click="goto(0)">全部文件</a>
      <template v-for="(v,i) of navs">
        <template v-if="i !== 0">
          <template v-if="i === navs.length-1">
            &nbsp;>&nbsp; {{v}}
          </template>
          <template v-else>
            &nbsp;>&nbsp; <a href="javascript:;" @click="goto(i)" :key=i>{{v}}</a>
          </template>
        </template>
      </template>
    </template>
  </div>

  <a-table 
  :columns="columns" 
  :data-source="resData.items"
  :pagination="false"
  :rowKey="(record,index) => index"
  :row-selection="{selectedRowKeys:selectedRowKeys,onChange:onSelectChange}"
  >
    <template slot="name" slot-scope="text, record,index">
      <a-icon type="folder" v-if="record.isDir"/>
      <a-icon type="file" v-else/>
      &nbsp;&nbsp;
      <template v-if="addFolder && index === 0">
          <a-input v-model="addFolderName" ref="addFolderInputRef" style="width:50%" @keyup.enter="handleAddFolderOK"/>&nbsp;
          <a-button icon="check" type="primary" size="small" @click="handleAddFolderOK"/>&nbsp;
          <a-button icon="close" type="primary" size="small" @click="handleAddFolderCancle"/>
      </template>
      <template v-else-if="renameIndex === index">
          <a-input v-model="renameValue" ref="renameInputRef" style="width:50%" @keyup.enter="handleRenameOK"/>&nbsp;
          <a-button icon="check" type="primary" size="small" @click="handleRenameOK"/>&nbsp;
          <a-button icon="close" type="primary" size="small" @click="handleRenameCancle"/>
      </template>
      <template v-else>
        <a v-if="record.isDir" href="javascript:;" @click="gonext(text)">{{text}}</a>
        <span v-else>{{text}}</span>
       </template>
    </template>
  </a-table>

  <!-- 上传列表 -->
  <div 
  id="upload_div"
  v-show="showTransfer">
    <div style="height: 48px;border-bottom: 2px solid #f6f6f6">
      <a-row justify="space-between" type="flex" style="line-height:40px">
        <a-col style="font-size:16px">上传列表 ( {{ this.uploadFinallyCount()}} / {{ Object.keys(this.uploadList).length }} )</a-col>
        <a-col style="margin-right:16px">
        <a-icon v-if="this.showTransferUploadList" type="minus" @click="()=>{ this.showTransferUploadList = false}"/>
        <a-icon v-else type="border" @click="()=>{ this.showTransferUploadList = true}"/>
        &nbsp;&nbsp;<a-icon type="close" @click="()=>{ this.showTransfer = false}"/>
        </a-col>
      </a-row>
    </div>
    <div 
    style="height:300px;overflow-y:auto;margin-top:5px"
    v-show="showTransferUploadList">
      <template v-for="v of this.uploadList">
        <a-row :key="v.key" class="upload_div_row">
          <a-col :span="6" style="text-overflow:ellipsis;overflow:hidden;white-space:nowrap">{{ v.filename }}</a-col>
          <a-col :span="17">
            <a-progress v-if="v.upSize < v.total" 
            :percent="v.upSize / v.total * 100" status="active"  />
            <a-progress v-else :percent="v.upSize / v.total * 100" status="success"  />
          </a-col>
        </a-row>
      </template>
    </div>
  </div>
  
  <!-- 移动对话框 -->
  <a-modal 
  v-model="dirModalVisible" 
  :title="dirModalTitle" 
  width="700px"
  centered
  :ok-text="dirModalOkText" 
  cancel-text="取消" 
  @ok="handleMvcp">
    <div class="path">
      <template v-if="dirModalNavs.length === 1">
        全部文件
      </template>
      <template v-else>
        <a href="javascript:;" @click="mvcpGoback()">返回上一级</a>
        <a-divider type="vertical" />
        <a href="javascript:;" @click="mvcpGoto(0)">全部文件</a>
        <template v-for="(v,i) of dirModalNavs">
          <template v-if="i !== 0">
            <template v-if="i === dirModalNavs.length-1">
              &nbsp;>&nbsp; {{v}}
            </template>
            <template v-else>
              &nbsp;>&nbsp; <a href="javascript:;" @click="mvcpGoto(i)" :key=i>{{v}}</a>
            </template>
          </template>
        </template>
      </template>
    </div>
    <div style="height:300px;overflow-y:auto">
      <a-table 
      :columns="dirModalColumns" 
      :data-source="dirModalData"
      :pagination="false"
      :showHeader="false"
      :rowKey="(record,index) => index"
      >
      <template slot="name" slot-scope="text">
        <a-icon type="folder" />
        &nbsp;&nbsp;
          <a  href="javascript:;" @click="mvcpGonext(text)">{{text}}</a>
      </template>
      </a-table>
    </div>
  </a-modal>

  <!-- 删除对话框 -->
  <a-modal 
  v-model="removeModalVisible" 
  title="确认删除" 
  width="450px"
  centered
  footer=""
  >
    <p style="text-align: center"><a-icon type="exclamation-circle" theme="twoTone" style="font-size:46px"/></p>
    <p style="text-align: center">确定删除所选文件吗？</p>
    <p style="text-align: center">
    <a-button shape="round" @click="()=>{this.removeModalVisible = false}">取消</a-button>&nbsp;
    <a-button type="primary" shape="round" @click="handleRemove">删除</a-button>
    </p>
  </a-modal>

  <!-- 分享对话框 -->
  <a-modal 
  v-model="shareModalVisible" 
  :title="shareTitle" 
  width="500px"
  centered
  footer=""
  >
    <div style="height:180px">
    <template v-if="shareCreate">
      <a-row>
        <a-col :span="4" style="height:54px;line-height:40px;">有效期：</a-col>
        <a-col :span="18"><a-slider v-model="shareTime" :marks="{1:'1天',7:'7天',30:'30天'}" :min="1" :max="30" /></a-col>
      </a-row>
      <br/><br/><br/>
      <p style="text-align: center"><a-button type="primary" shape="round" @click="handleShareCreate">创建链接</a-button></p>
    </template>
    <template v-else>
      <p style="text-align: center"> <a-input style="width: 300px"  v-model="sharedRoute" disabled></a-input></p>
      <p style="text-align: center">提取码 <a-input style="width: 60px" v-model="sharedToken" disabled></a-input></p>
      <br/>
      <p style="text-align: center">
        <a-button 
        type="primary" 
        shape="round" 
        v-clipboard:copy="sharedCopyText" 
        v-clipboard:success="()=>{this.sharedCopied = true}"
        >复制链接及提取码</a-button>
      </p>
      <p v-show="this.sharedCopied" style="text-align: center;color:#66B3FF">复制链接成功</p>
    </template>
    </div>
  </a-modal>
 </div>
</template>
<script>

import {mkdir,
list,
uploadCheck,
uploadFile,
fileRemove,
rename,
downloadFile,
sharedCreate,
mvcp} from "@/api/api"
import SparkMD5 from "spark-md5";

export default {
  name: 'filecoud',
  data () {
    return {
      columns:[
        {title:'文件名',dataIndex: 'filename',width:'60%',scopedSlots: { customRender: 'name' }},
        {title:'大小',dataIndex: 'size'},
        {title:'修改时间',dataIndex: 'date'},
      ],
      selectedRowKeys:[],
      selectedNames:[],

      root:"cloud",
      navs :[],
      resData:{},

      // 删除
      removeModalVisible: false,

      // 新建目录
      addFolder: false,
      addFolderName:"",

      // 上传文件
      uploadDirectory:false,
      uploadFileOpen:false,

      // 重命名
      renameIndex:-1,
      renameValue:"",

      // 移动、拷贝
      dirModalVisible:false,
      dirModalTitle:"",
      dirModalOkText:"",
      dirModalNavs:[],
      dirModalColumns:[
        {title:'文件名',dataIndex: 'filename',scopedSlots: { customRender: 'name' }},
      ],
      dirModalData:[],
      mvcpMove:false,
      mvcpSource:[],
      mvcpTarget:"",
      
      // 传输列表
      showTransfer:false,
      showTransferUploadList:true,
      uploadList:{},

      // 分享
      shareModalVisible:false,
      shareTitle:"",
      shareCreate:true,
      shareTime:1,
      sharedRoute:"",
      sharedToken:"",
      sharedDeadline:0,
      sharedCopyText:"",
      sharedCopied:false,
    }
  },
  mounted () {
    this.navs.push(this.root);
    this.goto(0);
  },
  methods:{
    onSelectChange (selectedRowKeys, selectedRows)  {
      //console.log(selectedRowKeys, selectedRows);
      this.selectedRowKeys = selectedRowKeys
      let names = [];
      for (let v of selectedRows){
        names.push(v.filename)
      }
      this.selectedNames = names;
    },
    getList(path){
      list({path:path}).then(res => {
        this.addFolder = false;
        this.addFolderName = "";
        this.selectedRowKeys = []
        this.selectedNames = []
        this.resData = res;
        // 排序 目录 > 名字
        this.resData.items.sort((a,b) =>{
          if ( !a.isDir && b.isDir){
            return 1
          }else if (a.isDir == b.isDir){
            return a.filename > b.filename ? 1 : -1
          }
          return -1
        })
       //console.log(res);
      })
    },
    goback(){
      if (this.navs.length > 1) {
        this.navs.pop()
        const path = this.navs.join("/")
        this.getList(path)
      }
    },
    gonext(dir){
      if (dir !== ""){
        this.navs.push(dir)
        const path = this.navs.join("/")
        this.getList(path)
      }
    },
    goto(index){
      if (index >= 0 && index < this.navs.length ){
        this.navs.splice(index+1)
        const path = this.navs.join("/")
        this.getList(path)
      }
    },
    progressColor (percent) {
      if (percent >= 80) {
        return 'red'
      } else if (percent >= 50) {
        return '#EAC100'
      }
    },
    openAddFolder(){
      if (!this.addFolder){
        const newRow = {
          filename:"",
          isDir:true,
          date:"",
          size:"-",
        }
        this.addFolder = true
        this.addFolderName = ""
        this.resData.items = [newRow,...this.resData.items]
        this.$nextTick(()=>{
          this.$refs.addFolderInputRef.focus()
        })
      }
    },
    handleAddFolderOK(){
      if (this.addFolder && this.addFolderName != ""){
        const path = this.navs.join("/") + "/" + this.addFolderName
        const args = { path:path}
        mkdir(args).then(() =>{
          const path = this.navs.join("/")
          this.getList(path)
        })
      }
    },
    handleAddFolderCancle(){
      if (this.addFolder){
        this.addFolder = false
        this.addFolderName = ""
        this.resData.items = this.resData.items.slice(1)
      }
    },
    handleUploadType(item){
      if (item.key === "folder"){
        this.uploadDirectory = true
      }else{
        this.uploadDirectory = false
      }
      this.uploadFileOpen = true
      this.$nextTick(() =>{
        this.$refs.upload.$refs.uploadRef.$el.firstChild.click();
        this.uploadFileOpen = false
      })
    },
    handleBeforeUpload(file){
      // console.log(file);
      // check

      const filename = file.name
      let path = this.navs.join("/")
      if (file.webkitRelativePath !== ""){
        const lastIndex = file.webkitRelativePath.lastIndexOf("/")
        if ( lastIndex > 0){
          path += "/" + file.webkitRelativePath.slice(0,lastIndex)
        }
      }

      // 分片大小 4M
      const sliceSize = 4 * 1024 *1024
      const size = file.size
      const total = Math.ceil(size / sliceSize)

      // md5
      this.getFileMd5(file,(fileMd5) => {
        const args = {path:path,filename:filename,md5:fileMd5,size:size,sliceTotal:total,sliceSize:sliceSize}
        uploadCheck(args).then(ret =>{
          this.showTransfer = true;
          this.showTransferUploadList = true;
          this.uploadList[fileMd5] = {
            key:fileMd5,
            filename:filename,
            total:size,
            upSize:0
          }
          //console.log(ret);
          // 开始上传
          if (ret.need) {
            const token = ret.token
            for (let current = 0;current < total ;current++){
              const index = current.toString()
              if (ret.existSlice === null || !(index in ret.existSlice)){
                let start = sliceSize * current,
                    end = start + sliceSize > file.size ? file.size : start + sliceSize;
                const blob = file.slice(start,end)

                let fd = new FormData();
                fd.append('path',path);
                fd.append('file', blob);
                fd.append('filename', filename);
                fd.append('current', index);
                fd.append('token', token);

                uploadFile(fd).then(() => {
                  // console.log(index,"ok");
                  this.updateProgress(fileMd5, end-start)
                })
              }else{
                this.updateProgress(fileMd5,sliceSize)
              }
            }
          }else{
            this.updateProgress(fileMd5,size)
          }
        })

      })
      //console.log(filename,path,size,total);
      return false
    },
    getFileMd5(file,callback){
      let fileSpark = new SparkMD5(),
          fileReader = new FileReader();

      fileReader.readAsBinaryString(file);
      fileReader.onload = function (ev) {
        fileSpark.appendBinary(ev.target.result);
        callback(fileSpark.end());
      };
      
    },
    updateProgress(key,value){
      const upInfo = this.uploadList[key]
      upInfo.upSize += value
      if (upInfo && upInfo.upSize >= upInfo.total){
          //console.log(upInfo);
          this.getList(this.navs.join("/"))
      }

      // 页面数据强制刷新
      this.$forceUpdate();
    },
    uploadFinallyCount(){
      let cnt = 0
      for (let k in this.uploadList){
        if (this.uploadList[k].upSize >= this.uploadList[k].total){
          cnt++
        }
      }
      return cnt
    },
    handleChange(info) {
      console.log(info);
      const status = info.file.status;
      if (status !== 'uploading') {
        console.log(info.file, info.fileList);
      }
      if (status === 'done') {
        this.$message.success(`${info.file.name} file uploaded successfully.`);
      } else if (status === 'error') {
        this.$message.error(`${info.file.name} file upload failed.`);
      }
    },
    openRemoveModal(){
      this.removeModalVisible = true
    },
    handleRemove(){
      let path = this.navs.join("/")
      fileRemove({path:path,filename:this.selectedNames}).then(()=>{
        this.getList(path)
      })
      .finally(()=>{
        this.removeModalVisible = false
      })
    },
    openRename(){
      if (this.selectedRowKeys.length === 1){
        if (this.addFolder){
          this.handleAddFolderCancle()
          this.selectedRowKeys = [this.selectedRowKeys[0]-1]
        }
        this.renameIndex = this.selectedRowKeys[0]
        this.renameValue = this.selectedNames[0]
        this.$nextTick(()=>{
          this.$refs.renameInputRef.focus()
        })
      }
      
    },
    handleRenameOK(){
      let path = this.navs.join("/")
      rename({path:path,oldName:this.selectedNames[0],newName:this.renameValue})
      .then(()=>{
        this.getList(path)
      })
      .finally(() => {
        this.handleRenameCancle()
      })
    },
    handleRenameCancle(){
      this.renameIndex = -1
    },
    openMvcpModal(move){
      if (this.selectedNames.length == 0){
        return
      }
      
      if (move){
        this.mvcpMove = true
        this.dirModalTitle="移动到"
        this.dirModalOkText="移动到此"
      }else{
        this.mvcpMove = false
        this.dirModalTitle="复制到"
        this.dirModalOkText="复制到此"
      }
      
      this.mvcpSource = []
      const path = this.navs.join("/")
      for (let name of this.selectedNames){
        this.mvcpSource.push(path + "/" + name)
      }

      this.dirModalNavs = [this.root]
      this.mvcpGetList(this.dirModalNavs.join("/"))
      //console.log(this.mvcpSource,this.dirModalVisible);
      this.dirModalVisible = true
    },
    mvcpGetList(path){
      list({path:path}).then(res => {
        this.dirModalData = [];
        for (let v of res.items){
          if (v.isDir){
            this.dirModalData.push(v)
          }
        }
        
        this.dirModalData.sort((a,b) =>{
          if ( !a.isDir && b.isDir){
            return 1
          }else if (a.isDir == b.isDir){
            return a.filename - b.filename
          }
          return -1
        })
        //console.log(this.dirModalData);
      })
    },
    mvcpGoto(index){
      if (index >= 0 && index < this.dirModalNavs.length ){
        this.dirModalNavs.splice(index+1)
        const path = this.dirModalNavs.join("/")
        this.mvcpGetList(path)
      }
    },
    mvcpGoback(){
      if (this.dirModalNavs.length > 1) {
        this.dirModalNavs.pop()
        const path = this.dirModalNavs.join("/")
        this.mvcpGetList(path)
      }
    },
    mvcpGonext(dir){
      if (dir !== ""){
        this.dirModalNavs.push(dir)
        const path = this.dirModalNavs.join("/")
        this.mvcpGetList(path)
      }
    },
    handleMvcp(){
      const args = {move:this.mvcpMove,source:this.mvcpSource,target:this.dirModalNavs.join("/")}
      mvcp(args)
      .then(()=>{
        this.getList(this.navs.join("/"))
      })
      .finally(()=>{
        this.dirModalVisible = false
      })
    },
    showDownload(){
      if (this.selectedRowKeys.length > 0){
        for (var idx of this.selectedRowKeys){
          if (this.resData.items[idx].isDir){
            return false
          }
        }
        return true
      }
      return false
    },
    handleDownload(){
      if (this.showDownload()){
        for (var idx in this.selectedNames){
          const filename = this.selectedNames[idx]
          const args = {path:this.navs.join("/"),filename:filename}
          downloadFile(args).then(res =>{
            //console.log(res);
            const blob = new Blob([res.data]);

            var downloadElement = document.createElement("a");
            var href = window.URL.createObjectURL(blob);
            downloadElement.href = href;
            downloadElement.download = decodeURIComponent(filename);
            document.body.appendChild(downloadElement);
            downloadElement.click();
            document.body.removeChild(downloadElement);
            window.URL.revokeObjectURL(href); 
          })
        }
      }
    },
    openShareModal(){
      if (this.selectedNames.length === 0){
        return
      }
      this.shareModalVisible = true;
      this.shareCreate = true;
      this.shareTime = 7
      this.shareTitle = "分享文件(夹):"
      this.sharedCopied = false
      if (this.selectedNames.length === 1){
        this.shareTitle += this.selectedNames[0]
      }else{
        this.shareTitle += this.selectedNames[0] + "等"
      }
    },
    handleShareCreate(){
      const args = {path:this.navs.join("/"),filename:this.selectedNames,deadline:this.shareTime}
      sharedCreate(args)
      .then((ret)=>{
        this.shareCreate = false;
        this.sharedRoute = ret.route;
        this.sharedToken = ret.sharedToken;
        this.sharedDeadline = ret.deadline;
        this.sharedCopyText = "链接：" + this.sharedRoute + "  提取码：" + this.sharedToken
      })
      
    }
  },

}
</script>
<style>
.body{
  padding:0 20px;
}
#header{
  height:60px;
  line-height:60px;
  margin-bottom:10px;
}
.path{
  margin-bottom:10px;
}

#upload_div{
  box-shadow: 0 1px 20px rgba(0,0,0,.2);
  border:2px;
  border-radius:4px;
  padding: 5px 16px;
  width:580px;
  position: fixed;
  bottom: 0;
  right: 20px;
  z-index: 1000;
  background: white;
}

.upload_div_row{
  height:40px;
  line-height:40px;
  background:#F0F0F0;
  padding:0 5px;
  margin-bottom:5px;
  border-radius:5px;
}

/* 滚动条宽度 */
::-webkit-scrollbar {
 width: 7px;
 height: 10px;
}
/* 滚动条的滑块 */
::-webkit-scrollbar-thumb {
 background-color: #a1a3a9;
 border-radius: 3px;
}
</style>