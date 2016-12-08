require! {
  'fs-extra' : fs
  'mkdirp'
  'path'
  'rimraf'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    rimraf.sync app-dir
    mkdirp.sync app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name),
                 app-dir



module.exports = ->
  @World = World
