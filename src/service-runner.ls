require! {
  'events' : {EventEmitter}
  'js-yaml' : yaml
  'fs'
  './next-port'
  'observable-process' : ObservableProcess
  'path'
}


class ServiceRunner extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(process.cwd!, @name, 'config.yml'), 'utf8')


  start: (done) ~>
    next-port (port) ~>
      @config['exorelay-port'] = port
      new ObservableProcess(@_create-start-command(@service-config.startup.command)
                            cwd: path.join(process.cwd!, @name),
                            verbose: yes,
                            console: log: @_log, error: @_log)
        ..wait @service-config.startup['online-text'], ~>
          @emit 'online', @name
          done!


  _create-start-command: (template) ->
    for key, value of @config
      template = template.replace "{{#{key}}}", value
    template


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceRunner
