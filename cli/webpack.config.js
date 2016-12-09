var path = require('path');
var fs = require('fs');

var nodeModules = {};
fs.readdirSync('node_modules')
  .filter(function(x) { return x !== '.bin'; })
  .forEach(function(mod) { nodeModules[mod] = 'commonjs ' + mod; });

module.exports = {
  entry: ['./src/cli'],
  target: 'node',
  node: {
    __filename: true,
    __dirname: true
  },
  output: {
    path: path.join(__dirname, 'build'),
    filename: 'bundle.js'
  },
  externals: nodeModules,
  module: {
    loaders: [
      {
        test: /\.ls$/,
        loader: 'livescript',
        exclude: /node_modules/
      },
      {
        test: /\.json$/,
        loader: 'json'
      }
    ],
    noParse: [/aws-sdk/]
  },
  resolve: {
    extensions: ["", ".js", ".json", ".ls"]
  }
}
