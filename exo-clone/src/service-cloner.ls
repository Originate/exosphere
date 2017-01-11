require! {
  'chalk': {red, green}
  'events' : EventEmitter
  'fs'
  'observable-process' : ObservableProcess
  'path'
}

class ServiceCloner extends EventEmitter

  ({@role, @config, @logger}) ->


  start: (done) ~>
    new ObservableProcess(@_create-command('git clone')
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code > 0            =>  @logger.log role: @role, text: red "Service cloning failed"
        | not @_is-valid-service!  =>  @logger.log role: @role, text: red "#{@role} is an invalid service"; exit-code = 1
        | _                        =>  @logger.log role: @role, text: green "done"
        done null, exit-code


  _is-valid-service: ->
    try
      fs.access-sync path.join(@config.path, 'service.yml')
      true
    catch
      false


  _create-command: (command) ->
    [command, @config.origin, @config.path].join ' '


  write: (text) ~>
    @logger.log {name: @role, text: text.trim!.replace(/\.*$/, '')}



module.exports = ServiceCloner
