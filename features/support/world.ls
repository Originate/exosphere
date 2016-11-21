require! {
  'fs-extra' : fs
  'path'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(process.cwd!, 'node_modules' 'exosphere-shared' 'example-apps', app-name), app-dir



module.exports = ->
  @World = World
