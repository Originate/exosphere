require! {
  'events' : EventEmitter
  'fs'
  'observable-process' : ObservableProcess
  'path'
}

class ServiceCloner extends EventEmitter

  (@name, @config) ->


  start: (done) ~>
    new ObservableProcess(@_create-command('git clone')
                          cwd: @config.root,
                          console: log: @_log, error: @_log)
      ..on 'ended', (exit-code) ~>
        | exit-code > 0            =>  @emit 'service-clone-fail', @name
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
    [command, @config.origin].join ' '


  _is-local-service: (path) ->
    path.substr(0, 2) is './'


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceCloner
