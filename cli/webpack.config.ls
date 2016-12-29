require!{
  'path'
  'fs'
}

node-modules = {}

fs.readdir-sync('node_modules')
  .filter((x) -> x isnt '.bin')
  .for-each((mod) -> node-modules[mod] = 'commonjs ' + mod)

module.exports = {
  entry: ['./src/cli']
  target: 'node'
  node:
    __filename: yes
    __dirname: yes
  output:
    path: path.join __dirname, 'build'
    filename: 'bundle.js'
  externals: node-modules
  module:
    loaders:
      * test: /\.ls$/
        loader: \livescript
        exclude: /node_modules/
      * test: /\.json$/
        loader: \json
    no-parse: [/aws-sdk/]
  resolve:
    extensions: ["", ".js", ".json", ".ls"]
}
