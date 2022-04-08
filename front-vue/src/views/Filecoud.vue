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
            <a-button icon="download" >离线下载</a-button>
          </a-button-group>
        </template>
        <template v-else>
          <a-button-group >
            <a-button icon="download" >下载</a-button>
            <a-button icon="delete" @click="handleRemove">删除</a-button>
            <a-button icon="form" v-show="this.selectedNames.length === 1" @click="openRename">重命名</a-button>
            <a-button icon="copy" @click="openMvcpModal(false)">复制</a-button>
            <a-button icon="drag" @click="openMvcpModal(true)">移动</a-button>
          </a-button-group>
        </template>
      </a-col>
      <a-col style="padding-right:200px">
        <a-popconfirm >
          <template slot="title">
          <br/>
            <div style="width:200px">

              test
              <a-button>test</a-button>
            </div>
          </template>
          <template slot="footer"></template>
          <a-icon type="swap" :rotate="90" @click="handleTransfer"/> 
        </a-popconfirm>
        
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
  :data-source="table.items"
  :pagination="false"
  :rowKey="(record,index) => index"
  :row-selection="{selectedRowKeys:selectedRowKeys,onChange:onSelectChange}"
  >
    <template slot="name" slot-scope="text, record,index">
      <a-icon type="folder" v-if="record.isDir"/>
      <a-icon type="file" v-else/>
      &nbsp;&nbsp;
      <template v-if="addFolder && index === 0">
          <a-input v-model="addFolderName"  style="width:50%"/>&nbsp;
          <a-button icon="check" type="primary" size="small" @click="handleAddFolderOK"/>&nbsp;
          <a-button icon="close" type="primary" size="small" @click="handleAddFolderCancle"/>
      </template>
      <template v-else-if="renameIndex === index">
          <a-input v-model="renameValue"  style="width:50%"/>&nbsp;
          <a-button icon="check" type="primary" size="small" @click="handleRenameOK"/>&nbsp;
          <a-button icon="close" type="primary" size="small" @click="handleRenameCancle"/>
      </template>
      <template v-else>
        <a v-if="record.isDir" href="javascript:;" @click="gonext(text)">{{text}}</a>
        <span v-else>{{text}}</span>
       </template>
    </template>
  </a-table>

  <a-modal 
  v-model="dirModalVisible" 
  :title="dirModalTitle" 
  width="700px"
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
mvcp} from "@/api/api"
import SparkMD5 from "spark-md5";

export default {
  name: 'filecoud',
  data () {
    return {
      columns:[
        {title:'文件名',dataIndex: 'filename',width:'60%',scopedSlots: { customRender: 'name' }},
        {title:'修改时间',dataIndex: 'date'},
        {title:'大小',dataIndex: 'size'},
      ],
      selectedRowKeys:[],
      selectedNames:[],

      root:"cloud",
      navs :[],
      table:{},

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
        this.table = res;
        // 排序 目录 > 名字
        this.table.items.sort((a,b) =>{
          if ( !a.isDir && b.isDir){
            return 1
          }else if (a.isDir == b.isDir){
            return a.filename - b.filename
          }
          return -1
        })
        //console.log(this.table);
      })
    },
    goback(){
      if (this.navs.length > 1) {
        this.navs.pop()
        const path = this.navs.join("/")
        //console.log(path);
        this.getList(path)
      }
    },
    gonext(dir){
      if (dir !== ""){
        this.navs.push(dir)
        const path = this.navs.join("/")
        //console.log(path);
        this.getList(path)
      }
    },
    goto(index){
      if (index >= 0 && index < this.navs.length ){
        this.navs.splice(index+1)
        const path = this.navs.join("/")
        //console.log(path);
        this.getList(path)
      }
    },
    openAddFolder(){
      const newRow = {
        filename:"",
        isDir:true,
        date:"",
        size:"-",
      }
      this.addFolder = true
      this.addFolderName = ""
      this.table.items = [newRow,...this.table.items]
    },
    handleAddFolderOK(){
      //console.log(this.addFolder,this.addFolderName);
      if (this.addFolder && this.addFolderName != ""){
        const path = this.navs.join("/") + "/" + this.addFolderName
        //console.log(path);
        const args = { path:path}
        //console.log(args);
        mkdir(args).then(() =>{
          //console.log(res);
          const path = this.navs.join("/")
          this.getList(path)
        })
      }
    },
    handleAddFolderCancle(){
      if (this.addFolder){
        this.addFolder = false
        this.addFolderName = ""
        this.table.items = this.table.items.slice(1)
      }
    },
    handleUploadType(item){
      if (item.key === "folder"){
        this.uploadDirectory = true
      }else{
        this.uploadDirectory = false
      }
      this.uploadFileOpen = true
      //console.log(item,this.uploadDirectory);
      this.$nextTick(() =>{
        this.$refs.upload.$refs.uploadRef.$el.firstChild.click();
        this.uploadFileOpen = false
      })
    },
    handleBeforeUpload(file){
      console.log(file);
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
          //console.log(ret);
          // 开始上传
          if (ret.need) {
            const token = ret.token
            let upCount = total
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
                  console.log(index,"ok");
                  upCount--
                  if (upCount == 0){
                    this.getList(this.navs.join("/"))
                  }
                })
              }else{
                upCount--
              }
            }
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
    handleTransfer(){

    },
    handleRemove(){
      let path = this.navs.join("/")
      fileRemove({path:path,filename:this.selectedNames}).then(()=>{
        this.getList(path)
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
      }
      
    },
    handleRenameOK(){
      let path = this.navs.join("/")
      rename({path:path,oldName:this.selectedNames[0],newName:this.renameValue})
      .then(()=>{
        this.getList(path)
      })
      .finally(() => {
        this.renameIndex = -1
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
}
.path{
  margin:10px 10px;
}
</style>