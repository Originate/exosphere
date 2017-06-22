require! {
  '../../../exosphere-shared' : {example-apps-path, run-process}
  'fs-extra' : fs
  'path'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(example-apps-path, app-name), app-dir


  @setup-app = (app-name, done) ->
    @run 'exo setup', path.join(process.cwd!, 'tmp', app-name)
      ..on 'ended', done


  @run = (command, app-dir) ->
    @process = run-process path.join(process.cwd!, 'bin', command), app-dir


module.exports = World
