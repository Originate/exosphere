require! {
  'child_process'
  'dim-console'
  '../../../exosphere-shared' : {templates-path}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
  'tmplconv'
}


observableProcessOptions = if process.env.DEBUG_EXOSPHERE_EXO_ADD
  stdout: dim-console.process.stdout
  stderr: dim-console.process.stderr
else
  stdout: no
  stderr: no


World = !->

  @create-empty-app = (app-name, done) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
      'exocom-version': '0.21.7'
    src-path = path.join templates-path, 'create-app'
    tmplconv.render(src-path, @app-dir, {data}).then ~> done!

  @run = (command) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: @app-dir,
                                     stdout: observableProcessOptions.stdout
                                     stderr: observableProcessOptions.stderr)


module.exports = ->
  @World = World
