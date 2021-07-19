<template>
  <div class="main">
    <h1>Blockchain simulator v0.0.0</h1>
    <el-form :inline="true" :model="newNetForm" class="demo-form-inline">
      <el-form-item label="网络名称:">
        <el-input v-model="newNetForm.name" placeholder="network name"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="success" @click="newNetwork">新建网络</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click=getNetworks type="primary">刷新</el-button>
      </el-form-item>
    </el-form>
    <el-row>
      <el-col :span="16" :offset="4">
      <el-table
        :data="networks"
        v-loading="loading"
        >
        <el-table-column
          prop="name"
          label="网络名称"
          align="center"
          >
        </el-table-column>
        <el-table-column
          prop="nodeNumber"
          label="节点数"
          align="center"
          >
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
            <router-link :to="'/network/'+scope.row.name">
              <el-button type="info" size="small">查看</el-button>
            </router-link>
            <el-button @click="delNetwork(scope.row.name)" type="danger" size="small">删除</el-button>
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
export default {
  name: 'Main',
  data () {
    return {
      networks:[],
      newNetForm:{
        name:'',
      },
      loading:false,
    }
  },
  created:function(){
    this.getNetworks()
  },
  methods:{
    getNetworks:function(){
      this.loading = true
      axios.get('/networks')
      .then(res=>{
        this.networks=res.data
      }).catch(err=>{
      }).finally(()=>{
        this.loading = false
      })
    },
    newNetwork:function(){
      if (this.newNetForm.name.length == 0){
        this.$message.error('网络名称不能为空')
        return
      }
      this.loading = true
      axios.get('/newnetwork',{params:{name:this.newNetForm.name}})
      .then(res=>{
        this.getNetworks()
      }).catch(err=>{
        if(err.response){
          this.$message.error('网络名称不能重复')
        }else{
          this.$message.error('网络连接失败')
        }
      }).finally(()=>{
        this.loading=false
      })
    },
    delNetwork:function(name){
      this.loading = true
      axios.get('/delnetwork',{params:{name:name}})
      .then(res=>{
        this.getNetworks()
      }).catch(err=>{
        this.$message.error('删除失败')
      }).finally(()=>{
        this.loading=false
      })
    },
  }
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
