require! {
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}



World = !->

  @checkout-app = (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name), @app-dir


  @setup-app = (@app-dir, done) ->
    new ObservableProcess(path.join(process.cwd!, '..' 'exo-setup' 'bin', 'exo-setup'),
                          cwd: @app-dir
                          stdout: dim-console.process.stdout
                          stderr: dim-console.process.stderr)
      ..on 'ended', done



module.exports = ->
  @World = World
