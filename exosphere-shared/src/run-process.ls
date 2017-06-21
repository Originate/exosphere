require! {
  'dim-console'
  './call-args'
  'observable-process': ObservableProcess
}

observableProcessOptions = if process.env.DEBUG_EXOSPHERE
  stdout: dim-console.process.stdout
  stderr: dim-console.process.stderr
else
  stdout: no
  stderr: no


module.exports = (command, dir) ->
  new ObservableProcess(call-args(command),
                        cwd: dir,
                        stdout: process.stdout
                        stderr: process.stderr)
