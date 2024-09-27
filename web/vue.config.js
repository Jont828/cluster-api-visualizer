// vue.config.js

const path = require('path')

module.exports = {
  css: {
    loaderOptions: {
      scss: {
        additionalData: `@use "@/styles/test.scss" as *;`,
      }
    }
  },
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
  },
}
