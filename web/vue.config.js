const { defineConfig } = require('@vue/cli-service')
const port = process.env.port || process.env.npm_config_port || 80 // 端口
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  devServer: {
  	host: '0.0.0.0',
	port: port,
  }
})
