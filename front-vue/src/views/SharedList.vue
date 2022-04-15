<template>
  <div class="body">
    <div id="header">
      <template v-if="this.selectedRowKeys.length === 0">个人云盘分享链接</template>
      <template v-else-if="this.selectedRowKeys.length === 1">
        <a-button type="primary" shape="round" icon="share-alt" 
          v-clipboard:copy="sharedCopyText" 
          v-clipboard:success="()=>{this.$message.success(`复制链接成功`);}"
        >复制链接</a-button>
        &nbsp;
        <a-button type="primary" shape="round" icon="share-alt" @click="handleCancel">取消分享</a-button>
      </template>
      <template v-else>
        <a-button type="primary" shape="round" icon="share-alt" @click="handleCancel">取消分享</a-button>
      </template>
    </div>

    <div class="path">全部文件</div>

    <a-table 
    :columns="columns" 
    :data-source="tableData"
    :pagination="false"
    :rowKey="(record) => record.key"
    :row-selection="{selectedRowKeys:selectedRowKeys,onChange:onSelectChange}"
    >
      <template slot="name" slot-scope="text, record">
        <a-icon type="appstore"/>&nbsp;&nbsp;{{showName(record)}}
      </template>
      <template slot="status" slot-scope="text, record">
        <span>{{showDeadline(record)}}后过期</span> 
      </template>
    </a-table>
  </div>
</template>

<script>
import moment from 'moment'
import {sharedList,sharedCancel} from "@/api/api"

export default {
  name: 'sharedlist',
  data () {
    return {
      columns:[
        {title:'文件',width:'45%',scopedSlots: { customRender: 'name' }},
        {title:'分享时间',width:'20%',dataIndex: 'createTime', customRender: (text) => moment.unix(text).format('YYYY-MM-DD HH:mm:ss')},
        {title:'状态', width:'15%',scopedSlots: { customRender: 'status' }},
        {title:'浏览次数',dataIndex: 'looked'},
      ],
      selectedRowKeys:[],
      tableData:[],
      sharedCopyText:'',
    }
  },
  mounted () {
    this.getList();
  },
  methods: {
    onSelectChange (selectedRowKeys, selectedRows)  {
      //console.log(selectedRowKeys, selectedRows);
      this.selectedRowKeys = selectedRowKeys
      if (this.selectedRowKeys.length === 1){
        // 分享链接
        const record = selectedRows[0]
        this.sharedCopyText = "链接：" + record.route + "  提取码：" + record.sharedToken
      }
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