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
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'config.yml'), 'utf8')


  start: (done) ~>
    port-reservation.get-port N (port) ~>
      @config['exorelay-port'] = port
      new ObservableProcess(@_create-command(@service-config.startup.command)
                            cwd: @config.root,
                            verbose: yes,
                            console: log: @_log, error: @_log)
        ..wait @service-config.startup['online-text'], ~>
          @emit 'online', @name
          done!


  _create-command: (template) ->
    template = "#{@config.root}/#{template}"
    for key, value of @config
      template = template.replace "{{#{key}}}", value
    template


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceRunner
