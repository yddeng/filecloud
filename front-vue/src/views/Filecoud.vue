<template>
 <div class="body">
  <div id="header">
    <a-button type="primary" icon="upload">上传文件</a-button>&nbsp;
    <a-button icon="folder-add">新建文件夹</a-button>&nbsp;
    <a-button icon="download">下载文件</a-button>
  </div>

  <div id="path">
    <template v-if="navs.length === 1">
      全部文件
    </template>
    <template v-else>
      <a href="#" @click="goback()">返回上一级</a>
      <a-divider type="vertical" />
      <a href="#" @click="goto(0)">全部文件</a>
      <template v-for="(v,i) of navs">
        <template v-if="i !== 0">
          <template v-if="i === navs.length-1">
            &nbsp;>&nbsp; {{v}}
          </template>
          <template v-else>
            &nbsp;>&nbsp; <a href="#" @click="goto(i)" :key=i>{{v}}</a>
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
  >
    <template slot="name" slot-scope="text, record">
      <a href="#" @click="gonext(text)" v-if="record.isDir">{{text}}</a>
      <span v-else>{{text}}</span>
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
      navs :["cloud"],
      table:{}
    }
  },
  mounted () {
    this.mkdir("cloud/nav1/nav2");
    this.goto(0);
  },
  methods:{
    mkdir (path ){
      const args = { path:path}
      console.log(args);
      mkdir(args)
      .then(res =>{
        console.log(res);
      })
      .catch(err => console.log(err))
    },
    getList(path){
      list({path:path})
      .then(res => {
         console.log(res);
         this.table = res.data;
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
    }
  },

}
</script>
<style>
.body{
  padding:0 10px;
}
#header{
  height:40px;
  line-height:40px;
}
#path{
  margin:10px 0;
}
</style>