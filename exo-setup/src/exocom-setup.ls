require! {
  'child_process'
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper}
  'observable-process' : ObservableProcess
}


class ExocomSetup extends EventEmitter

  (@logger) ->
    @name = \exocom


  start: ~>
    version = child_process.exec-sync 'npm show exocom-dev version' |> (.to-string!) |> (.trim!)
    if DockerHelper.image-exists author: \originate, name: \exocom, version: version
      @logger.log role: @name, text: 'ExoCom image already up to date'
      return
    @logger.log role: @name, text: "Pulling ExoCom image version #{version}"
    new ObservableProcess((DockerHelper.get-pull-command author: \originate, name: \exocom, version: version),
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code is 0  =>  @logger.log role: @name, text: "ExoCom image updated to version #{version}"
        | otherwise       =>  @logger.log role: @name, text: "Failed to retrieve latest ExoCom image"


  write: (text) ~>
    @logger.log {@name, text, trim: yes}



module.exports = ExocomSetup
