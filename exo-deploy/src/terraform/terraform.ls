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


  apply: ({hosted-zone-id}, done) ->
    var-flags = if hosted-zone-id then "-var 'hosted_zone_id=#{hosted-zone-id}'" else ''
    new ObservableProcess("terraform apply #{var-flags}",
                          cwd: '/usr/src/terraform')
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform apply failed: #{exit-code}")
        done!


  destroy: (done) ->
    # remove -Xnew-destroy flag when hashicorp/terraform docker image is updated to 0.8
    new ObservableProcess("terraform destroy -force -Xnew-destroy",
                          cwd: '/usr/src/terraform')
      ..wait 'Enter a value:', ~>
        ..enter 'yes'
      ..on 'ended', (exit-code) ->
        | exit-code  =>  return done new Error("terraform destroy failed: #{exit-code}")
        done!


module.exports = Terraform
