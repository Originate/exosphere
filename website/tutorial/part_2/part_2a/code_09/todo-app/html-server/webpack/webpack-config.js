const path = require('path');

module.exports = {

  entry: {
    main: path.resolve(__dirname, '../app/client/index.js')
  },

  output: {
    path: path.resolve(__dirname, '../public/assets'),
    publicPath: '/assets/',
    filename: '[name].js',
    pathinfo: true
  },

  resolve: {
    extensions: ['.js', '.styl']
  },

  module: {
    loaders: [
      {
        test: /\.styl/,
        use: [
          { loader: 'style-loader' },
          { loader: 'css-loader' },
          { loader: 'stylus-loader'},
        ]
      }
    ]
  }

}
