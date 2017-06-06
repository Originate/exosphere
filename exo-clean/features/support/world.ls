require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'wait': {wait}
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


  @add-file = (@app-dir, @service-name, @file-name, done) ->
    app-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'application.yml'), 'utf8')
    service-config = app-config.services[\public][@service-name] or app-config.services[\private][@service-name]
    fs.write-file-sync path.join(@app-dir, service-config.location, @file-name), 'test'
    wait 1_000 done


module.exports = ->
  @World = World
