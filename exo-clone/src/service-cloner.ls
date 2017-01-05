require! {
  'chalk': {red, green}
  'events' : EventEmitter
  'fs'
  'observable-process' : ObservableProcess
  'path'
}

class ServiceCloner extends EventEmitter

  ({@name, @config, @logger}) ->


  start: (done) ~>
    new ObservableProcess(@_create-command('git clone')
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code > 0            =>  @logger.log name: @name, text: red "Service cloning failed"
        | not @_is-valid-service!  =>  @logger.log name: @name, text: red "#{@name} is an invalid service"; exit-code = 1
        | _                        =>  @logger.log name: @name, text: green "#{@name} cloned successfully"
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
    @logger.log {@name, text: text.trim!.replace(/\.*$/, '')}



module.exports = ServiceCloner
