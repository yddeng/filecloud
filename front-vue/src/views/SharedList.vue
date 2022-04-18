<template>
  <div class="body">
    <div id="header">
      <template v-if="this.selectedRowKeys.length === 0">个人云盘分享链接</template>
      <template v-else-if="this.selectedRowKeys.length === 1">
        <a-button type="primary" shape="round" icon="paper-clip" 
          v-clipboard:copy="sharedCopyText" 
          v-clipboard:success="handleCopyOk"
        >复制链接</a-button>
        &nbsp;
        <a-button type="primary" shape="round" icon="stop" @click="handleCancel">取消分享</a-button>
      </template>
      <template v-else>
        <a-button type="primary" shape="round" icon="stop" @click="handleCancel">取消分享</a-button>
      </template>
    </div>

    <div class="path">全部文件</div>

    <a-table 
    :columns="columns" 
    :data-source="tableData"
    :pagination="false"
    :rowKey="(record) => record.key"
    :row-selection="{selectedRowKeys:selectedRowKeys,onChange:onSelectChange}"
    :customRow="customRowFunc"
    >
      <template slot="name" slot-scope="text, record">
        <a-icon type="appstore"/>&nbsp;&nbsp;{{showName(record)}}
      </template>
      <template slot="status" slot-scope="text, record">
        <span v-if="record.key !== hoverKey">{{showDeadline(record)}}后过期</span> 
        <div v-else>
          <a 
            v-clipboard:copy="hoverCopyText" 
            v-clipboard:success="handleCopyOk"
          ><a-icon type="paper-clip" />&nbsp;复制链接</a>&nbsp;&nbsp;&nbsp;
          <a @click="handleCancelKey"><a-icon type="stop" />&nbsp;取消分享</a>
        </div> 
      </template>
      <template slot="looked" slot-scope="text, record">
        <span v-if="record.key !== hoverKey">{{text}}</span> 
      </template>
    </a-table>
  </div>
</template>

<script>
import moment from 'moment'
import {sharedList,sharedCancel} from "@/api/api"
import {makeSharedLink} from "@/utils/common"

export default {
  name: 'sharedlist',
  data () {
    return {
      columns:[
        {title:'文件',width:'45%',scopedSlots: { customRender: 'name' }},
        {title:'分享时间',width:'20%',dataIndex: 'createTime', customRender: (text) => moment.unix(text).format('YYYY-MM-DD HH:mm:ss')},
        {title:'状态', width:'15%',scopedSlots: { customRender: 'status' }},
        {title:'浏览次数',dataIndex: 'looked',scopedSlots: { customRender: 'looked' }},
      ],
      selectedRowKeys:[],
      tableData:[],
      hoverKey:'',
      hoverCopyText:'',
      sharedCopyText:'',
    }
  },
  mounted () {
    this.getList();
  },
  methods: {
    onSelectChange (selectedRowKeys, selectedRows)  {
      console.log(selectedRowKeys, selectedRows);
      this.selectedRowKeys = selectedRowKeys
      this.sharedCopyTextUpdate()
    },
    sharedCopyTextUpdate(){
      if (this.selectedRowKeys.length === 1){
        // 分享链接
        for (let record of this.tableData){
          if (record.key === this.selectedRowKeys[0]){
            this.sharedCopyText =  makeSharedLink(record.route, record.sharedToken)
            break
          }
        }
      }
    },
    customRowFunc(record){
      return {
        props: {},
        on: { // 事件
          click: (event) => {
            //console.log(event);
            if (event.target.localName === 'a'){
              // 阻止事件
              return
            }
            let selectedIdx = this.selectedRowKeys.indexOf(record.key)
            if ( selectedIdx === -1 ){
              this.selectedRowKeys.push(record.key)
            }else{
              this.selectedRowKeys.splice(selectedIdx,1)
            }
            //console.log( selectedIdx ,this.selectedRowKeys,selectedIdx );
            this.sharedCopyTextUpdate()
          },       // 点击行
          dblclick: () => {},
          contextmenu: () => {},
          mouseenter: () => {
            this.hoverKey = record.key
            this.hoverCopyText = makeSharedLink(record.route, record.sharedToken)
          },  // 鼠标移入行
          mouseleave: () => {this.hoverKey = ''}
        },

      };
    },
    showName(record){
      let name = record.filename[0];
      if (record.filename.length > 1){
        name += "等"
      }
      return name
    },
    showDeadline(record){
      let start = moment.unix(record.createTime)
      let end = moment.unix(record.deadline)
      return end.to(start,true)
    },
    getList () {
      sharedList({}).then(res => {
        this.selectedRowKeys = []
        this.tableData = []
        for (let v in res){
          this.tableData.push(res[v]);
        }
        this.tableData.sort((a,b) =>{
          return a.createTime - b.createTime
        })
        //console.log(this.tableData);
      })
    },
    handleCancel(){
      if (this.selectedRowKeys.length == 0){
        return
      }

      sharedCancel({keys:this.selectedRowKeys})
      .then(()=>{
        this.$message.success(`分享已取消`);
        this.getList()
      })
    },
    handleCopyOk(){
      this.$message.success(`复制链接成功`);        
    },
    handleCancelKey(){
      if (this.hoverKey === ''){
        return
      }

      sharedCancel({keys:[this.hoverKey]})
      .then(()=>{
        this.$message.success(`分享已取消`);
        this.getList()
      })
    }
  }
}
</script>

<style>
.body{
  padding:10px 30px;
}
#header{
  height:40px;
  line-height:40px;
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