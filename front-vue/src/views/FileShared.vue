<template>
  <div class="body">
  
    <div id="header">
      <div id="header-title">
        <a-row justify="space-between" type="flex">
          <a-col style="font-size:20px;">
            <a-icon type="appstore" theme="twoTone" /> &nbsp; <span>{{this.sharedInfo.name}}</span>
          </a-col>
          <a-col  style="margin-right:40px">
            <a-button icon="download" :disabled="!showDownload()" @click="handleDownload">下载</a-button>
          </a-col>
        </a-row>
        
      </div>
      <div id="header-time">
        <a-icon type="clock-circle" /> <span>{{showCreateTime()}}</span> &nbsp;&nbsp; <span>过期时间：{{showDeadline()}}后</span> 
      </div>
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
    :data-source="tableData"
    :pagination="false"
    :rowKey="(record,index) => index"
    :row-selection="{selectedRowKeys:selectedRowKeys,onChange:onSelectChange}"
    >
      <template slot="name" slot-scope="text, record">
        <a-icon type="folder" v-if="record.isDir"/>
        <a-icon type="file" v-else/>
        &nbsp;&nbsp;
          <a v-if="record.isDir" href="javascript:;" @click="gonext(text)">{{text}}</a>
          <span v-else>{{text}}</span>
      </template>
    </a-table>
  

    <!-- token -->
    <a-modal 
    v-model="sharedModalToken" 
    title="提取分享链接" 
    width="500px"
    centered
    :maskStyle="{opacity:1,background:'#F0F0F0'}"
    :maskClosable="false"
    :closable="false"
    footer=""
    >
      <div style="height:180px">
        <br/>
        <p>请输入提取码：</p>
        <a-input v-model="sharedModalValue" style="width:300px" @keyup.enter="handleSharedInfo"></a-input>&nbsp;
        <a-button type="primary" @click="handleSharedInfo">提取文件</a-button>
        <p v-show="this.sharedTokenCheckFailed" style="color:	#FF5151;font-size:12px">提取码错误</p>

      </div>
    </a-modal>
  </div>
</template>

<script>
import {
sharedInfo,
sharedList,
sharedDownload,
} from "@/api/api"
import moment from 'moment'

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

      sharedToken:"",
      sharedKey:"",
      sharedInfo:{},

      navs :[],
      tableData:[],

      sharedModalToken:false,
      sharedModalValue:"",
      sharedTokenCheckFailed:false,
    }
  },
  mounted () {
    console.log(this.$route);
    this.sharedKey = "5Pa2WpBcXeYqz3lf"
    if (this.sharedToken === ""){
      this.sharedModalToken = true
    }
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
      sharedList({path:path,key:this.sharedKey,sharedToken:this.sharedToken}).then(res => {
        this.selectedRowKeys = []
        this.selectedNames = []
        this.tableData = res.items;
        // 排序 目录 > 名字
        this.tableData.sort((a,b) =>{
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
    showDeadline(){
      let start = moment.unix(this.sharedInfo.createTime)
      let end = moment.unix(this.sharedInfo.deadline)
      return end.to(start,true)
    },
    showCreateTime(){
      return moment.unix(this.sharedInfo.createTime).format('YYYY-MM-DD HH:mm:ss')
    },
    handleSharedInfo(){
      const args = {key:this.sharedKey,sharedToken:this.sharedModalValue}
      sharedInfo(args)
      .then((ret) => {
        //console.log(ret);
        this.sharedModalToken = false
        this.sharedToken = args.sharedToken
        this.sharedInfo = ret
        this.tableData = ret.items
        this.navs.push(ret.root)
      })
      .catch(()=>{
        this.sharedTokenCheckFailed = true
      })
    },
    showDownload(){
      if (this.selectedRowKeys.length > 0){
        for (var idx of this.selectedRowKeys){
          if (this.tableData[idx].isDir){
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
          const args = {key:this.sharedKey,sharedToken:this.sharedToken,path:this.navs.join("/"),filename:filename}
          sharedDownload(args).then(res =>{
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
  },
}
</script>
<style>
.body{
  padding:10px 30px;
}
#header{
  height:86px;
  margin-bottom:10px;
  border-bottom: 1px solid #F0F0F0;
}

#header-title{
  height:60px;
  line-height:60px;
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