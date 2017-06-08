require! {
  '../../../exosphere-shared' : {run-process}
  'child_process'
  'path'
}


World = !->

  @make-repo = (cwd) ->
    child_process.exec-sync("git init && git add --all && git commit -m \"initial commit\"",
                            cwd: cwd)


  @run = (command) ->
    @process = run-process path.join(process.cwd!, 'bin', command),
                           path.join(process.cwd!, 'tmp')


module.exports = ->
  @World = World
