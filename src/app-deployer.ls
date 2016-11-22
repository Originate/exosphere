require! {
  'events' : {EventEmitter}
  'observable-process' : ObservableProcess
  'path'
  'tmplconv'
}


# Deploys the overall application
class AppDeployer extends EventEmitter

  (@app-config) ->


  generate-terraform: ->


  deploy: ->
    new ObservableProcess("terraform apply",
                          cwd: path.join process.cwd!, 'terraform'
                          stdout: process.stdout,
                          stderr: process.stderr)

module.exports = AppDeployer
