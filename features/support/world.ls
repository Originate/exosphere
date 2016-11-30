require! {
  'dim-console'
  'exosphere-shared' : {call-args, example-apps-path}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(example-apps-path, app-name), app-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', 'exo setup'),
                                     cwd: path.join(process.cwd!, 'tmp', app-name),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



module.exports = ->
  @World = World
