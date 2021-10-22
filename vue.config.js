// vue.config.js
module.exports = {
  configureWebpack: {
    module: {
      rules: [{ test: /\.yaml$/, use: 'raw-loader' }],
    },
  },

  transpileDependencies: [
    'vuetify'
  ]
}
