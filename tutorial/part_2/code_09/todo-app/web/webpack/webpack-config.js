const path = require('path');

module.exports = {

  entry: {
    main: path.resolve('app/client')
  },

  output: {
    path: path.resolve(__dirname, '../public/assets'),
    publicPath: '/assets/',
    filename: '[name].js',
    pathinfo: true
  },

  resolve: {
    extensions: ['', '.js', '.ls', '.styl']
  },

  module: {
    loaders: [
      {
        test: /\.ls$/,
        loader: 'livescript'
      }, {
        test: /\.styl/,
        loader: 'style!css!stylus'
      }
    ]
  }

}
