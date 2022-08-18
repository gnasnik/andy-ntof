<template>
    <div id="app-container">
    <el-row style="text-align:left;margin: 0 0 20px;">
      <el-card style="padding: 20px 20px" shadow="never">
        <el-col :xs="12" :sm="16" :md="8" :lg="8" :xl="8">
          <label class="label">参与人数</label>
          <div class="value">
            <span>{{stat.Players}}</span>
        </div>
        </el-col>
        <el-col :xs="12" :sm="16" :md="8" :lg="8" :xl="8">
          <label class="label">商品数量</label>
          <div class="value">
            <span>{{stat.GoodCount}}</span>
          </div>

        </el-col>
        <el-col :xs="12" :sm="16" :md="8" :lg="8" :xl="8">
          <label class="label">市值</label>
          <div class="value"> {{stat.MarketCap.toFixed(2)}} </div>
        </el-col>
      </el-card>
    </el-row>        
        <el-card>
            <linechart :chartData="chartData"></linechart>
        </el-card>
        <el-card style="margin-top: 20px;">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="所有" name="all"></el-tab-pane>
            <el-tab-pane label="玩家" name="user"></el-tab-pane>
          </el-tabs>
            <el-input v-model="name" prefix-icon="el-icon-search" size="small" placeholder="查询玩家名称" v-if="activeName==='user'" @change="getListByUser"></el-input>

            <el-table
                :data="dataList"
                style="width: 100%">
                <el-table-column
                    prop="Date"
                    label="日期"
                    >
                </el-table-column>
                <el-table-column
                    prop="Name"
                    label="姓名"
                    >
                  <template slot-scope="scope">
                    {{ scope.row.Name }}
                    <el-tag v-if="scope.row.isVIP" type="danger" size="mini"> VIP </el-tag>
                  </template>
                </el-table-column>
                <el-table-column
                    prop="Count"
                    label="数量"
                   >
                </el-table-column>                
                <el-table-column
                    prop="Asset"
                    label="投入">
                </el-table-column>
                </el-table>
        </el-card>
    </div>

</template>

<script>
import linechart from './linechart.vue'
import { getStats, getPlayers, getUserList } from '@/api/api'
export default {
  components: { linechart },
  data() {
      return {
          name: 'iisutas',
          activeName: 'all',
          stat: undefined,
          dataList:[],
          chartData: {
              label: [],
              data: {
                  players: [],
                  cap: []
              },
          },
          vips: ["芝士", "SXY77","悠悠", "我是陈伟霆","bigdick86", "gk", "水冰月", "小丸子", "旺仔牛奶", "Tiandix", "陈胖子", "Katrina",
          "呵呵呀呵呵", "Lns", "TAing", "甜甜","windy", "伟业", "尼古拉斯·广坤", "尼古拉斯·赵四", "一西", "妮9", "木森","Hong", "shaos", 
	        "黄不了"],
      }
  },
  created(){
    this.getList()
    this.stats()
  },
  methods:{
      handleClick(tab) {
        if (tab.name === 'user') {
           this.getListByUser()
        }else {
          this.getList()
        }

      },
      stats(){
          getStats({}).then(res =>{
              if (res) {
                  this.chartData.label = res.data.label
                  this.chartData.data.players = res.data.players
                  this.chartData.data.cap = res.data.cap
                  this.stat = res.data.stats[res.data.stats.length-1]
              }
          })
      },
      getList() {
          this.dataList = []
          getPlayers().then(res => {
              if (res) {
                  for (let i = 0; i< res.data.players.length; i++) {
                      let p = res.data.players[i]
                      if (this.vips.includes(p.Name)) {
                         p.isVIP = true
                      }
                      this.dataList.push(p)
                  }
              }
          })
      },
      getListByUser(){
        if (this.name === '') {
          return
        }

        getUserList({name:this.name}).then( res => {
          if (res) {
            this.dataList = res.data.list
          }
        })
      }
  }
    
}
</script>

<style lang="css" scoped>
  .label {
    color: #3f475a;
    font-size: 18px;
    letter-spacing: -.5px;
    line-height: 24px;
    font-weight: lighter;
  }

  .value {
    color: #1f2533;
    font-size: 20px;
    font-weight: 500;
    line-height: 40px;
  }
</style>
