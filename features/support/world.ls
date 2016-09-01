require! {
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(process.cwd!, 'node_modules' 'exosphere-shared' 'example-apps', app-name), app-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'node_modules' 'exo-setup' 'bin' 'exo-setup'),
                                     cwd: path.join(process.cwd!, 'tmp', app-name),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



module.exports = ->
  @World = World
