require! {
  'dim-console'
  'observable-process' : ObservableProcess
  'path'
  'prelude-ls' : {any}
  'wait' : {wait-until}
}


observableProcessOptions = if process.env.DEBUG_EXOSERVICE
  verbose: yes
  stdout: process.stdout
  stderr: process.stderr
else
  verbose: no
  stdout: no
  stderr: no


# Provides steps for end-to-end testing of the service as a stand-alone binary
CliWorld = !->

  @create-exoservice-instance = ({role, exocom-port}, done) ->
    command = "#{process.cwd!}/bin/exo-js"
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(command,
                                     env: {EXOCOM_PORT: exocom-port, ROLE: role},
                                     cwd: path.join(process.cwd!, 'features', 'example-services', role),
                                     verbose: observableProcessOptions.verbose,
                                     stdout: observableProcessOptions.stdout,
                                     stderr: observableProcessOptions.stderr)
      ..wait 'online at port', done


  @remove-register-service-message = (@exocom, done) ->
    wait-until (~> @exocom.received-messages.length), 10, ~>
      @exocom.reset! if @exocom.received-messages |> any (.name is 'exocom.register-service')
      done!

module.exports = CliWorld
