// vue.config.js
module.exports = {
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = 'Cluster API Visualizer'
        return args
      })
  },
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      "^/api": {
        target: "http://0.0.0.0:8081",
        changeOrigin: true,
        logLevel: 'debug' 
      },
    },
  },
}
