<template>
  <div class="network">
    <el-row>
      <el-col :offset="4" :span="8">
        <el-page-header @back="back" :content="'Network: '+$route.params.name">
        </el-page-header>
      </el-col>
    </el-row>
    <el-form :inline="true" :model="newNodeForm" class="demo-form-inline">
      <el-form-item label="节点名称:">
        <el-input v-model="newNodeForm.name" placeholder="node name"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="success" @click="newNode">新增节点</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="getNodes" type="primary">刷新</el-button>
      </el-form-item>
    </el-form>

    <h1>网络拓扑</h1>
    <div style="width: 100%; height: 60vh" ref="chart"></div>

    <el-row>
      <el-col :span="16" :offset="4">
      <el-table
        :data="nodes"
        v-loading="loading"
        >
        <el-table-column
          prop="id"
          label="id"
          align="center"
          >
        </el-table-column>
        <el-table-column
          prop="name"
          label="节点名称"
          align="center"
          >
        </el-table-column>
        <el-table-column
          prop="connNumber"
          label="连接数"
          align="center"
          >
        </el-table-column>
        <el-table-column
          prop="up"
          label="状态"
          align="center"
          >
          <template slot-scope="scope">
            <el-tag v-if="scope.row.up" type="success">运行</el-tag>
            <el-tag v-else type="danger">停止</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="nodeType"
          label="节点类型"
          align="center"
          >
          <template slot-scope="scope">
            <el-tag>{{scope.row.nodeType}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          align="center"
          >
          <template slot-scope="scope">
            <el-link @click="startNode(scope.row.id)" type="success" size="small">启动</el-link>
            <el-link @click="stopNode(scope.row.id)" type="warning" size="small">停止</el-link>
            <el-link @click="delNode(scope.row.id)" type="danger" size="small">删除</el-link>
          </template>
        </el-table-column>
      </el-table>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import axios from 'axios'
import {simnetPrefix} from '../api/simnet'
axios.defaults.baseURL = simnetPrefix

import * as echarts from 'echarts'

export default {
  name: 'Main',
  data () {
    return {
      network:"",
      nodes:[],
      loading:false,
      chart:null,
      newNodeForm:{
        name:'',
      },
    }
  },
  mounted(){
    this.network = this.$route.params.name
    this.createNetworkGraph()
    this.getNodes()
  },
  methods:{
    back: function(){
      this.$router.push('/')
    },
    getNodes(){
      this.loading = true
      axios.get("/nodes/"+this.network)
      .then(res=>{
        this.nodes=res.data
      }).catch(err=>{
      }).finally(()=>{
        this.loading = false
      })
    },
    newNode(){
      this.loading = true
      axios.post("/nodes/"+this.network,{
        name:this.newNodeForm.name,
      })
      .then(res=>{
        this.nodes.push(res.data)
      }).catch(err=>{
        if(err.response){
          this.$message.error(err.response.data)
        }else{
          this.$message.error("unknown error")
        }
      }).finally(()=>{
        this.loading = false
      })
    },
    startNode(id){
      this.loading = true
      axios.post("/nodes/"+this.network+"/"+id+"/start")
      .then(res=>{
        this.getNodes()
      }).catch(err=>{
      }).finally(()=>{
        this.loading = false
      })
    },
    stopNode(id){
      this.loading = true
      axios.post("/nodes/"+this.network+"/"+id+"/stop")
      .then(res=>{
        this.getNodes()
      }).catch(err=>{
      }).finally(()=>{
        this.loading = false
      })
    },
    delNode(id){
      this.loading = true
      axios.delete("/nodes/"+this.network+"/"+id)
      .then(res=>{
        this.getNodes()
      }).catch(err=>{
      }).finally(()=>{
        this.loading = false
      })
    },
    createNetworkGraph(){
      this.chart = echarts.init(this.$refs.chart);
      let myChart = this.chart

      let data = [{
          fixed: true,
          x: myChart.getWidth() / 2,
          y: myChart.getHeight() / 2,
          symbolSize: 20,
          id: '-1'
      }];

      let edges = [];

      let option = {
          series: [{
              type: 'graph',
              layout: 'force',
              animation: false,
              data: data,
              force: {
                  // initLayout: 'circular'
                  // gravity: 0
                  repulsion: 100,
                  edgeLength: 5
              },
              edges: edges
          }]
      };
      setInterval(function () {
          data.push({
              id: data.length
          });
          var source = Math.round((data.length - 1) * Math.random());
          var target = Math.round((data.length - 1) * Math.random());
          if (source !== target) {
              edges.push({
                  source: source,
                  target: target
              });
          }
          myChart.setOption({
              series: [{
                  roam: false,
                  data: data,
                  edges: edges
              }]
          });
      }, 500);
      myChart.setOption(option)
    },
  }
}
function randoms(length,chars){
    var maxNum=chars.length-1;
    var hex='';
    var num=0;
    for(let i=0;i<length;i++){
        num=rand(0,maxNum-1);
        hex+=chars.slice(num,num+1);  
    }
return hex;
}
function rand(minNum,maxNum){
    var choices=maxNum-minNum;
    let minMax = 0
    var num=minMax+Math.round(Math.random()*choices)
    return num;
}
function delArrVal(arr,val){    //查找数组中的某个值并全部删除    第一个参数是查找的数组 第二个参数是删除的值
  for(let i=0;i<arr.length;i++){
    if(arr[i].id==val){
      arr.splice(i,1)
      i--;
    } 
  } 
  return arr;
} 
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
