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
import { getStats, getPlayers } from '@/api/api'
export default {
  components: { linechart },
  data() {
      return {
          stat: undefined,
          dataList:[],
          chartData: {
              label: [],
              data: {
                  players: [],
                  cap: []
              },
          },
      }
  },
  created(){
    this.getList()
    this.stats()
  },
  methods:{
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
          getPlayers().then(res => {
              if (res) {
                  this.dataList = res.data.players
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