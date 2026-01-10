// vue.config.js
module.exports = {
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = 'Cluster API Visualizer'
        return args
      })
    config.module
      .rule("mjs")
      .test(/\.mjs$/)
      .type("javascript/auto")
      .include.add(/node_modules/)
      .end();
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
    client: {
      webSocketTransport: 'sockjs',
      webSocketURL: 'auto://0.0.0.0:0/ws',
      overlay: {
        errors: true,
        warnings: false,
      },
    },
    webSocketServer: 'sockjs',
    hot: true,
    liveReload: true,
  },
}
