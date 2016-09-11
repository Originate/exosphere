require! {
  'child_process'
  'dim-console'
  'exosphere-shared' : {templates-path}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmplconv'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(process.cwd!, 'node_modules' 'exosphere-shared' 'example-apps', app-name), app-dir


  @create-empty-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join templates-path, 'create-app'
    tmplconv.render(src-path, app-dir, {data})


  @create-repo = (repo-name) ->
    repos-dir = path.join process.cwd!, 'tmp', 'repos'
    fs.empty-dir-sync repos-dir
    fs.copy-sync path.join(process.cwd!, 'node_modules' 'exosphere-shared' 'example-apps/test', repo-name),
                 path.join(repos-dir, repo-name)
    child_process.exec-sync "git init #{repo-name}", cwd: repos-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'node_modules' 'exo-setup' 'bin', 'exo-setup'),
                                     cwd: path.join(process.cwd!, 'tmp', app-name),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



module.exports = ->
  @World = World
