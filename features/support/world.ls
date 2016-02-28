require! {
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
}


World = !->

  @checkout-app = (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, 'example-apps', app-name), @app-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(path.join('..', 'bin', 'exo-install'),
                                     cwd: path.join(process.cwd!, 'tmp'),
                                     verbose: yes,
                                     console: dim-console)
      ..wait "installation complete", done


  @start-app = (app-name, done) ->
    @process = new ObservableProcess(path.join('..', 'bin', 'exo-run'),
                                     cwd: path.join(process.cwd!, 'tmp'),
                                     verbose: yes,
                                     console: dim-console)
      ..wait "application ready", done



module.exports = ->
  @World = World
