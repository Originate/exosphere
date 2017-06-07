require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args, example-apps-path}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}

observableProcessOptions = if process.env.DEBUG_EXOSPHERE_CLI
  stdout: dim-console.process.stdout
  stderr: dim-console.process.stderr
else
  stdout: no
  stderr: no


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(example-apps-path, app-name), app-dir


  @setup-app = (app-name, done) ->
    @run 'exo setup', path.join(process.cwd!, 'tmp', app-name)
      ..on 'ended', done


  @run = (command, app-dir) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: observableProcessOptions.stdout
                                     stderr: observableProcessOptions.stderr)


module.exports = ->
  @World = World
