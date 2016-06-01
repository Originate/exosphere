require! {
  'events' : {EventEmitter}
  'js-yaml' : yaml
  'fs'
  'nitroglycerin' : N
  'port-reservation'
  'observable-process' : ObservableProcess
  'path'
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
        ..wait @service-config.startup['online-text'], ~>
          @emit 'online', @name
          done!


  _create-command: (template) ->
    if @_is_local_command template
      template = "#{@config.root}/#{template}"
    template

  _is_local_command: (command) ->
    command.substr(0, 2) is './'


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceRunner
