const HtmlWebpackPlugin = require('html-webpack-plugin')
const path = require('path')

module.exports = {
  entry: "./src/index.js",
  output: {
    filename: '[hash].js',
    path: path.resolve(__dirname, 'dist'),
    publicPath: '',
  },
  plugins: [
    new HtmlWebpackPlugin({
      backendUrl: process.env.BACKEND_EXTERNAL_ORIGIN,
      template: 'src/index.pug',
    })
  ],
  module: {
    rules: [{
      test: /\.pug$/,
      loader: 'pug-loader',
    }],
  },
}
