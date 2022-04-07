<template>
 <div class="body">
  <div id="header">
    <a-row justify="space-between" type="flex">
      <a-col>
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
        <a-button icon="folder-add" @click="handleAddFolder">新建文件夹</a-button>
        <a-button icon="download" >下载文件</a-button>
        </a-button-group>
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

  <div id="path">
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
  rowKey="(record,index) => index"
  :row-selection="rowSelection"
  >
    <template slot="name" slot-scope="text, record,index">
      <template v-if="record.isDir">
        <template v-if="addFolder && index === 0">
          <a-icon type="folder" />&nbsp;&nbsp;
          <a-input v-model="addFolderName"  style="width:50%"/>&nbsp;
          <a-button icon="check" type="primary" size="small" @click="handleAddFolderOK"/>&nbsp;
          <a-button icon="close" type="primary" size="small" @click="handleAddFolderCancle"/>
        </template>
        <template v-else>
          <a-icon type="folder" />&nbsp;&nbsp;<a href="javascript:;" @click="gonext(text)">{{text}}</a>
        </template>
      </template>
      <template v-else><a-icon type="file" />&nbsp;&nbsp;<span >{{text}}</span></template>
    </template>
  </a-table>
 </div>
</template>
<script>

import {mkdir,list} from "@/api/api"

export default {
  name: 'filecoud',
  data () {
    return {
      columns:[
        {title:'文件名',dataIndex: 'filename',width:'60%',scopedSlots: { customRender: 'name' }},
        {title:'修改时间',dataIndex: 'date'},
        {title:'大小',dataIndex: 'size'},
      ],
      rowSelection:{
        onChange: (selectedRowKeys, selectedRows) => {
          console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
        onSelect: (record, selected, selectedRows) => {
          console.log(record, selected, selectedRows);
        },
        onSelectAll: (selected, selectedRows, changeRows) => {
          console.log(selected, selectedRows, changeRows);
        },
      },
      navs :["cloud"],
      table:{},
      addFolder: false,
      addFolderName:"",
      uploadDirectory:false,
      uploadFileOpen:false,
    }
  },
  mounted () {
    this.mkdir("cloud/nav1/nav2");
    this.goto(0);
  },
  methods:{
    getList(path){
      list({path:path})
      .then(res => {
         console.log(res);
         this.addFolder = false;
         this.addFolderName = "";
         this.table = res;
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
    mkdir (path ){
      const args = { path:path}
      console.log(args);
      mkdir(args)
      .then(() =>{
        //console.log(res);
        const path = this.navs.join("/")
        this.getList(path)
      })
      .catch(err => console.log(err))
    },
    handleAddFolder(){
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
      console.log(this.addFolder,this.addFolderName);
      if (this.addFolder && this.addFolderName != ""){
        const path = this.navs.join("/") + "/" + this.addFolderName
        console.log(path);
        this.mkdir(path)
      }
    },
    handleAddFolderCancle(){
      this.addFolder = false
      this.addFolderName = ""
      this.table.items = this.table.items.slice(1)
    },
    handleUploadType(item){
      if (item.key === "folder"){
        this.uploadDirectory = true
      }else{
        this.uploadDirectory = false
      }
      this.uploadFileOpen = true
      console.log(item,this.uploadDirectory);
      this.$nextTick(() =>{
        this.$refs.upload.$refs.uploadRef.$el.firstChild.click();
        this.uploadFileOpen = false
      })
    },
    handleBeforeUpload(file,fileList){
      console.log(file,fileList);
    },
    handleChange(info) {
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

    }
  },

}
</script>
<style>
.body{
  padding:0 20px;
}
#header{
  height:50px;
  line-height:50px;
}
#path{
  margin:10px 10px;
}
</style>