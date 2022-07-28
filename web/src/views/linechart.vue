<template>
  <div :class="className" :style="{height:height,width:width}" />
</template>

<script>
  import * as echarts from 'echarts'
  require('echarts/theme/macarons') // echarts theme

  export default {
    props: {
      className: {
        type: String,
        default: 'chart'
      },
      width: {
        type: String,
        default: '100%'
      },
      height: {
        type: String,
        default: '350px'
      },
      autoResize: {
        type: Boolean,
        default: true
      },
      chartData: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        chart: null
      }
    },
    watch: {
      chartData: {
        deep: true,
        handler(val) {
          this.setOptions(val)
        }
      }
    },
    mounted() {
      this.$nextTick(() => {
        this.initChart()
      })
    },
    beforeDestroy() {
      if (!this.chart) {
        return
      }
      this.chart.dispose()
      this.chart = null
    },
    methods: {
      initChart() {
        this.chart = echarts.init(this.$el, 'macarons')
        this.setOptions(this.chartData)
      },
      setOptions({
        label,
        data
      } = {}) {
        // label = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
        // data = {
        //   a: [0, 24.1, 7.2, 15.5, 50.4, 53.4, 42.4, 24.1, 7.2, 15.5, 42.4, 55],
        //   b: [0, 41.1, 30.4, 65.1, 33.3, 43.3, 53.3, 41.1, 30.4, 65.1, 53.3, 70],

        // }

        this.chart.setOption({
          title: {
            text: '市值曲线',
            x: 'left',
            textStyle: {
              fontSize: '15',
              color: '#000'
            },
          },
          grid: {
            top: 70,
            right: 20,
            bottom: 30,
            left: 30
          },
          tooltip: {
            trigger: 'axis',
            backgroundColor: 'rgba(255, 255, 255)'
          },
          legend: {
            data: [ '总市值'],
            right: 0
          },
          xAxis: {
            data: label,
            axisLabel: {
              margin: 2,
              lineStyle: {
                color: '#7b8293',
              },
              textStyle: {
                color: '#7b8293',
                fontSize: 12
              }
            }
          },
          yAxis: [{
            type: 'value',
            // name: '价格',
            splitLine: {
              show: true,
              lineStyle: {
                type: 'dashed',
                color: '#f5f5f5'
              }
            },
            axisLabel: {
              margin: 2,
              formatter: function (value, index) {
                if (value >= 1000 && value < 1000000) {
                  return value / 1000 + 'k';
                } else if (value >= 1000000) {
                  return value / 1000000 + 'M';
                } else  {
                  return value
                }
              },
              lineStyle: {
                color: '#7b8293',
              },
              textStyle: {
                color: '#7b8293',
                fontSize: 12
              }
            }
          }, ],
          series: [
            {
              name: '总市值',
              type: 'line',
              symbolSize: 6,
              symbol: 'circle',
              smooth: true,
              data: data.cap,
              lineStyle: {
                color: '#205fec'
              },
              itemStyle: {
                color: '#205fec',
                borderColor: '#205fec'
              },
              areaStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                    offset: 0,
                    color: '#205fecb3'
                  },
                  {
                    offset: 1,
                    color: '#205fec03'
                  },
                ]),
              },
            //   emphasis: {
            //     itemStyle: {
            //       color: {
            //         type: 'radial',
            //         x: 0.5,
            //         y: 0.5,
            //         r: 0.5,
            //         colorStops: [{
            //             offset: 0,
            //             color: '#9E87FF'
            //           },
            //           {
            //             offset: 0.4,
            //             color: '#9E87FF'
            //           },
            //           {
            //             offset: 0.5,
            //             color: '#fff'
            //           },
            //           {
            //             offset: 0.7,
            //             color: '#fff'
            //           },
            //           {
            //             offset: 0.8,
            //             color: '#fff'
            //           },
            //           {
            //             offset: 1,
            //             color: '#fff'
            //           },
            //         ],
            //       },
            //       borderColor: '#9E87FF',
            //       borderWidth: 2,
            //     },
            //   },
            },
          ],
        })
      }
    }
  }
</script>