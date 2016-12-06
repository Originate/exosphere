require! {
  'observable-process' : ObservableProcess
  'path'
}

# Encapsulates running Terraform on the local machine
class Terraform


  get: ({target}, done) ->
    new ObservableProcess("terraform get",
                          cwd: target or '/usr/src/terraform')
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


  output: ({variable}, done) ->
    process = new ObservableProcess("terraform output #{variable}",
                                    stdout: no,
                                    cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done null, new Error("terraform output failed: #{exit-code}")
        done process.full-output!


  # executes 'terraform destroy' at target, ignoring the provided hosted zone
  destroy: ({target, hosted-zone-id}, done) ->
    var-flags = if hosted-zone-id then "-var 'hosted_zone_id=#{hosted-zone-id}'" else ''
    console.log "var flags: #{var-flags}"
    new ObservableProcess("terraform destroy -force #{var-flags}" + target,
                          cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform apply failed: #{exit-code}")
        done!


module.exports = Terraform
