require! {
  '../../../exosphere-shared' : {run-process}
  'js-yaml' : yaml
  'fs-extra' : fs
  'path'
  'wait': {wait}
}



World = !->

  @checkout-app = (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name), @app-dir


  @setup-app = (@app-dir, done) ->
    run-process path.join(process.cwd!, '..' 'exo-setup' 'bin' 'exo-setup')
      ..on 'ended', done


  @add-file = (@app-dir, @service-name, @file-name, done) ->
    app-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'application.yml'), 'utf8')
    service-config = app-config.services[\public][@service-name] or app-config.services[\private][@service-name]
    fs.write-file-sync path.join(@app-dir, service-config.location, @file-name), 'test'
    wait 1_000 done


module.exports = World
