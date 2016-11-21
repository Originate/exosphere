require! {
  'child_process'
  'events' : {EventEmitter}
  'observable-process' : ObservableProcess
}


class ExocomSetup extends EventEmitter

  (@logger) ->
    @name = \exocom


  start: ~>
    #TODO: Update to use docker-helper class from exosphere-shared
    version = child_process.exec-sync 'npm show exocom-dev version' |> (.to-string!) |> (.trim!)
    if child_process.exec-sync 'docker images originate/exocom' |> (.to-string!) |> (.includes version)
      @logger.log name: @name, text: 'ExoCom image already up to date'
      return
    @logger.log name: @name, text: "Pulling ExoCom image version #{version}"
    new ObservableProcess("docker pull originate/exocom:#{version}",
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>  @logger.log name: @name, text: "ExoCom image updated to version #{version}"
        | otherwise       =>  @logger.log name: @name, text: "Failed to retrieve latest ExoCom image"


  write: (text) ~>
    @emit 'output', {@name, text, trim: yes}



module.exports = ExocomSetup
