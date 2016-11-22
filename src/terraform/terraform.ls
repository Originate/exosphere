require! {
  'observable-process' : ObservableProcess
  'path'
}

# Encapsulates running Terraform on the local machine
class Terraform


  get: (done) ->
    new ObservableProcess("terraform get",
                          cwd: path.join(process.cwd!, 'terraform')
                          stdout: process.stdout,
                          stderr: process.stderr)
      ..on 'ended', (exit-code) ->
        | exit-code  =>  done new Error("terraform get failed: #{exit-code}")
        done!


  apply: ->
    new ObservableProcess("terraform apply",
                          cwd: path.join(process.cwd!, 'terraform')
                          stdout: process.stdout,
                          stderr: process.stderr)


module.exports = Terraform
