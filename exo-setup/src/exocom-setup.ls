require! {
  'chalk' : {red}
  'child_process'
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper}
  'observable-process' : ObservableProcess
}


class ExocomSetup extends EventEmitter

  (@app-config, @logger) ->
    @name = \exocom


  start: ~>
    version = @app-config.bus.version
    if DockerHelper.image-exists author: \originate, name: \exocom, version: version
      @logger.log role: @name, text: 'ExoCom image already up to date'
      return
    @logger.log role: @name, text: "Pulling ExoCom image version #{version}"
    new ObservableProcess((DockerHelper.get-pull-command author: \originate, name: \exocom, version: version),
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code is 0  =>  @logger.log role: @name, text: "ExoCom image updated to version #{version}"
        | otherwise       =>  throw new Error red "Failed to retrieve latest ExoCom image"


  write: (text) ~>
    @logger.log {role: @name, text, trim: yes}



module.exports = ExocomSetup
