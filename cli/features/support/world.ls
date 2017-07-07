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


  @checkout-service-template = (app-name, template-name) ->
    src = path.join(process.cwd!, '..' 'exosphere-shared' 'templates', 'boilr-templates', template-name)
    dest = path.join process.cwd!, 'tmp', app-name, '.exosphere', template-name
    fs.copy-sync src, dest


  @setup-app = (app-name, done) ->
    @run 'exo setup', path.join(process.cwd!, 'tmp', app-name)
      ..on 'ended', done


  @run = (command, app-dir) ->
    @process = run-process path.join(process.cwd!, 'bin', command), app-dir


module.exports = World
