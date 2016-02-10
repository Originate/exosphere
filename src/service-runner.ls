require! {
  'events' : {EventEmitter}
  'js-yaml' : yaml
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


class ServiceRunner extends EventEmitter

  (@name, @config) ->


  start: (done) ~>
    service-config = yaml.safeLoad fs.readFileSync(path.join(process.cwd!, @name, 'config.yml'), 'utf8')
    new ObservableProcess(@_create-start-command(service-config.startup.command)
                                     cwd: path.join(process.cwd!, @name),
                                     verbose: yes,
                                     console: log: @_log, error: @_log)
      ..wait service-config.startup['online-text'], ~>
        @emit 'online', @name
        done!


  _create-start-command: (template) ->
    for key, value of @config
      template = template.replace "{{#{key}}}", value
    template


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceRunner
