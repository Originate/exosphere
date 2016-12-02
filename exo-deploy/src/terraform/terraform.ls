require! {
  'observable-process' : ObservableProcess
  'path'
}

# Encapsulates running Terraform on the local machine
class Terraform


  get: (done) ->
    new ObservableProcess("terraform get",
                          cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform get failed: #{exit-code}")
        done!


  pull-remote-state: ({backend, backend-config}, done) ->
    options = "-backend=#{backend} "
    for config in backend-config
      options += "-backend-config=#{config} "
    new ObservableProcess("terraform remote config #{options}",
                          cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform remote config failed: #{exit-code}")
        done!


  apply: (done) ->
    new ObservableProcess("terraform apply",
                          cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform apply failed: #{exit-code}")
        done!


module.exports = Terraform
