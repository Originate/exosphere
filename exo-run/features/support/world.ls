require! {
  '../../../exosphere-shared' : {run-process}
  'fs-extra' : fs
  'path'
}



class World

  checkout-app: (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name), @app-dir


  setup-app: (@app-dir, done) ->
    run-process path.join(process.cwd!, '..' 'exo-setup' 'bin', 'exo-setup'), @app-dir
      ..on 'ended', done


  run-app: ({command, online-text}, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @app-dir
    if online-text then
      @process.wait online-text, done
    else
      done!


  checkout-and-run-app: ({online-text}, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      command = \exo-run
      if process.platform is 'win32' then command += '.cmd'
      @run-app {command, online-text}, done


module.exports = World
