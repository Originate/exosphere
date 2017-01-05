require! {
  'child_process'
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
  'tmplconv'
}

# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


World = !->

  @make-repo = (cwd) ->
    child_process.exec-sync("git init && git add --all && git commit -m \"initial commit\"",
                            cwd: cwd,
                            stdio: [1,2])



module.exports = ->
  @app-dir = app-dir
  @World = World
