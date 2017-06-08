require! {
  'child_process'
  '../../../exosphere-shared' : {example-apps-path, templates-path}
  'fs-extra' : fs
  'mkdirp'
  'path'
  'tmplconv'
}


World = !->

  @checkout-app = (@app-name) ->
    app-dir = path.join process.cwd!, 'tmp', @app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(example-apps-path, @app-name), app-dir


  @create-empty-app = (@app-name) ->
    app-dir = path.join process.cwd!, 'tmp', @app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': @app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join templates-path, 'create-app'
    tmplconv.render(src-path, app-dir, {data})


  @create-origin = (service-type, service-dir) ->
    origins-dir = path.join process.cwd!, 'tmp', 'origins'
    service-origin-dir = path.join origins-dir, service-type

    # create the origin directory
    mkdirp.sync origins-dir
    fs.rename-sync path.join(service-dir),
                   service-origin-dir
    child_process.exec-sync "git init", cwd: service-origin-dir

    # make the origin directory the new origin of the service
    child_process.exec-sync "git clone #{service-origin-dir}", cwd: @app-dir
    service-origin-dir


  @setup-app = (@app-name, done) ->
    command = path.join(process.cwd!, 'node_modules' 'exo-setup' 'bin', 'exo-setup')
    @process = run-process command, path.join(process.cwd!, 'tmp', @app-name)
      ..on 'ended', done



module.exports = ->
  @World = World
