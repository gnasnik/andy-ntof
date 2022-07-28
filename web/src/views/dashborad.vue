<template>
    <div id="app-container">
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