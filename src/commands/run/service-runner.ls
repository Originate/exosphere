require! {
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
}


class ServiceRunner extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')
    @service-config.messages.sends ?= []
    @service-config.messages.receives ?= []


  start: (done) ~>
    port-reservation.get-port N (port) ~>
      @config.SERVICE_NAME = @name
      @config.EXORELAY_PORT = port
      new ObservableProcess(@_create-command(@service-config.startup.command)
                            cwd: @config.root,
                            env: @config
                            verbose: yes,
                            console: log: @_log, error: @_log)
        ..on 'ended', ~> throw new Error "Service '#{@name}' crashed"
        ..wait @service-config.startup['online-text'], ~>
          @emit 'online', @name
          done!


  _create-command: (command) ->
    if @_is_local_command command
      command = path.join @config.root, command
    command


  _is_local_command: (command) ->
    command.substr(0, 2) is './'


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceRunner
