require!{
  'path'
  'fs'
}

node-modules = {}
apps = {}

fs.readdir-sync('node_modules')
  .filter((x) -> x isnt '.bin')
  .for-each((mod) -> node-modules[mod] = 'commonjs ' + mod)

module.exports = {
  entry:
    'exo-js': './src/exo-service.ls'
  target: 'node'
  node:
    __filename: yes
    __dirname: yes
  output:
    path: path.join process.cwd!, 'dist'
    filename: '[name].js'
    libraryTarget: 'commonjs2'
    library: 'ExoService'
  externals: node-modules
  module:
    rules:
      * test: /\.ls$/
        use: <[ livescript-loader ]>
        exclude: /node_modules/
      ...
  resolve:
    extensions: [".js", ".json", ".ls"]
}
