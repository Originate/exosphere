require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}



class World

  checkout-app: (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name), @app-dir


  setup-app: (@app-dir, done) ->
    new ObservableProcess(path.join(process.cwd!, '..' 'exo-setup' 'bin', 'exo-setup'),
                          cwd: @app-dir
                          stdout: dim-console.process.stdout
                          stderr: dim-console.process.stderr)
      ..on 'ended', done


  run-app: ({command, online-text}, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
    if online-text then
      @process.wait online-text, done
    else
      done!


  checkout-and-run-app: ({@app-name, online-text}, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      command = \exo-run
      if process.platform is 'win32' then command += '.cmd'
      @run-app {command, online-text}, done


module.exports = ->
  @World = World
