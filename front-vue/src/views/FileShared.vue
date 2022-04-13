<template>
  <div class="body">
    <div id="header">
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
  </div>
</template>

<script>
import {
list,
} from "@/api/api"
export default {
  name: 'fileshared',
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