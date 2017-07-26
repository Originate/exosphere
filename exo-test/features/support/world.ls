require! {
  'fs-extra' : fs
  '../../../exosphere-shared' : {run-process}
  'path'
}


World = !->

  @checkout-app = (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name), @app-dir


  @setup-app = (app-name, done) ->
    command = "exo run"
    @process = run-process command, path.join(process.cwd!, 'tmp', app-name)
      ..wait 'setup complete', done


module.exports = World
