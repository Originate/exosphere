require! {
  'chalk': {red}
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
        | exit-code > 0            =>  @logger.log @name, text: red "Service cloning failed"
        | not @_is-valid-service!  =>  @emit 'service-invalid', @name; exit-code = 1
        | _                        =>  @emit 'service-clone-success', @name
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
    @logger.log {@name, text: (text.trim!.replace /\.*$/, '')}



module.exports = ServiceCloner
