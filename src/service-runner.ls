require! {
  'events' : {EventEmitter}
  'exosphere-shared' : {call-args}
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
      new ObservableProcess(call-args(@_create-command @service-config.startup.command),
                            cwd: @config.root,
                            env: @config
                            stdout: {@write} 
                            stderr: {@write})
        ..on 'ended', ~> throw new Error "Service '#{@name}' crashed"
        ..wait @service-config.startup['online-text'], ~>
          @emit 'online', @name
          done!


  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @config.root, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  write: (text) ~>
    @emit 'output', {@name, text}


module.exports = ServiceRunner
