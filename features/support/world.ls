require! {
  'child_process'
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(process.cwd!, 'example-apps', app-name), app-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', 'exo setup'),
                                     cwd: path.join(process.cwd!, 'tmp', app-name),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @start-app = (app-name, done) ->
    @process = new ObservableProcess(path.join('..', 'bin', 'exo run'),
                                     cwd: path.join(process.cwd!, 'tmp'),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..wait "application ready", done


  @make-repo = (cwd) ->
    child_process.exec-sync("git init && git add --all && git commit -m \"initial commit\"",
                            cwd: cwd,
                            stdio: [1,2])



module.exports = ->
  @World = World
