<template>
  <div class="body">
    <div id="header">
      <template v-if="this.selectedRowKeys.length === 0">个人云盘分享链接</template>
      <template v-else-if="this.selectedRowKeys.length === 1">
        <a-button type="primary" shape="round" icon="share-alt" >复制链接</a-button>
        &nbsp;
        <a-button type="primary" shape="round" icon="share-alt" @click="handleCancel">取消分享</a-button>
      </template>
      <template v-else>
        <a-button type="primary" shape="round" icon="share-alt" @click="handleCancel">取消分享</a-button>
      </template>
    </div>

    <div class="path">
        全部文件
    </div>

    <a-table 
    :columns="columns" 
    :data-source="tableData"
    :pagination="false"
    :rowKey="(record) => record.key"
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
  </div>
</template>

<script>
import {sharedList,sharedCancel} from "@/api/api"

export default {
  name: 'sharedlist',
  data () {
    return {
      columns:[
        {title:'文件',width:'50%'},
        {title:'分享时间',dataIndex: 'createTime'},
        {title:'状态', },
        {title:'浏览次数',dataIndex: 'looked'},
      ],
      selectedRowKeys:[],
      tableData:[],
    }
  },
  mounted () {
    this.getList();
  },
  methods: {
    onSelectChange (selectedRowKeys)  {
      //console.log(selectedRowKeys, selectedRows);
      this.selectedRowKeys = selectedRowKeys
    },
    getList () {
      sharedList({}).then(res => {
        this.selectedRowKeys = []
        this.tableData = []
        //console.log(res);
        for (let v in res){
          this.tableData.push(res[v]);
        }
        
        this.tableData.sort((a,b) =>{
          return a.createTime - b.createTime
        })
        console.log(this.tableData);
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