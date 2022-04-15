<template>
  <div class="body">
    <div v-if="!sharedChecked"></div>
    <div v-else>
    <div id="header">
      <div id="header-title">
        <template v-if="!showDownload()">
          <a-icon type="appstore" theme="twoTone" style="font-size:20px;"/> &nbsp; <span style="font-size:20px;">{{this.sharedInfo.name}}</span>
        </template>
        <template v-else>
          <a-button type="primary" icon="download" shape="round" @click="handleDownload">下载</a-button>
        </template> 
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
    :customRow="customRowFunc"
    >
      <template slot="name" slot-scope="text, record">
        <a-icon type="folder" v-if="record.isDir"/>
        <a-icon type="file" v-else/>
        &nbsp;&nbsp;
          <a v-if="record.isDir" href="javascript:;" @click="gonext(text)">{{text}}</a>
          <span v-else>{{text}}</span>
      </template>
    </a-table>
    </div>

    <!-- token -->
    <a-modal 
    v-model="sharedModalToken" 
    width="500px"
    centered
    :maskStyle="{opacity:1,background:'#FCFCFC'}"
    :maskClosable="false"
    :closable="false"
    footer=""
    >
      <div style="height:240px">
        <div style="height:50px;line-height:36px;text-align:center;border-bottom: 2px solid #F0F0F0;">
          <img src="@/assets/logo.png" alt="logo" style="width:28px;height:28px;">&nbsp;
          <span style="font-size:24px;font-weight: 600;">个人云盘</span>
        </div>
        <p style="font-size:16px;margin-top:40px">请输入提取码：</p>
        <a-input v-model="sharedModalValue" style="width:320px" size="large" @keyup.enter="handleSharedInfo"></a-input>&nbsp;
        <a-button type="primary" size="large" @click="handleSharedInfo">提取文件</a-button>
        <p v-show="this.sharedTokenCheckFailed" style="color:	#FF5151;font-size:12px">{{sharedTokenCheckText}}</p>

      </div>
    </a-modal>
  </div>
</template>

<script>
import {
sharedInfo,
sharedPath,
sharedDownload,
} from "@/api/api"
import moment from 'moment'
import  storage from 'store'

export default {
  name: 'fileshared',
  data () {
    return {
      columns:[
        {title:'文件名',dataIndex: 'filename',width:'50%',scopedSlots: { customRender: 'name' }},
        {title:'大小',dataIndex: 'size',width:'20%',},
        {title:'修改时间',dataIndex: 'date'},
      ],
      selectedRowKeys:[],
      selectedNames:[],

      sharedChecked:false,
      sharedToken:"",
      sharedKey:"",
      sharedInfo:{},

      navs :[],
      tableData:[],

      sharedModalToken:false,
      sharedModalValue:"",
      sharedTokenCheckFailed:false,
      sharedTokenCheckText:"",
    }
  },
  mounted () {
    //console.log(this.$route);
    this.sharedKey = this.$route.params.key
    this.sharedToken = storage.get(this.sharedKey)
    
    if (this.sharedToken){
      const args = {key:this.sharedKey,sharedToken:this.sharedToken}
      sharedInfo(args)
      .then((ret) => {
        this.sharedInfo = ret
        this.tableData = ret.items
        this.navs.push(ret.root)
        this.sharedChecked = true
      })
      .catch(()=>{
        storage.remove(this.sharedKey)
        this.sharedModalToken = true
      })
    }else{
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
    customRowFunc(record,index){
      return {
        props: {},
        on: { // 事件
          click: () => {
            //console.log(event);
            // if (event.target.localName === 'a'){
            //   // 阻止事件
            //   return
            // }
            let selectedIdx = this.selectedRowKeys.indexOf(index)
            if ( selectedIdx === -1 ){
              this.selectedRowKeys.push(index)
              this.selectedNames.push(record.filename)
            }else{
              this.selectedRowKeys.splice(selectedIdx,1)
              this.selectedNames.splice(selectedIdx,1)
            }
            //console.log(selectedIdx ,this.selectedRowKeys );
            
          },       // 点击行
          dblclick: () => {},
          contextmenu: () => {},
          mouseenter: () => {},  // 鼠标移入行
          mouseleave: () => {}
        },

      };
    },
    getList(path){
      sharedPath({path:path,key:this.sharedKey,sharedToken:this.sharedToken}).then(res => {
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
        storage.set(this.sharedKey,this.sharedToken,12 * 60 *60 * 1000)
        this.sharedInfo = ret
        this.tableData = ret.items
        this.navs.push(ret.root)
        this.sharedChecked = true
      })
      .catch((err)=>{
        this.sharedTokenCheckFailed = true
        this.sharedTokenCheckText = err.data.message
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