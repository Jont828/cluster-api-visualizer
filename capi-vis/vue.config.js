// vue.config.js
module.exports = {
  configureWebpack: {
    module: {
      rules: [{ test: /\.yaml$/, use: 'raw-loader' }],
    },
  },
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api': {
        target: 'http://localhost:3080',
        changeOrigin: true
      },
    }
  }
}
