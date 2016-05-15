require! {
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


class ServiceInstaller extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    @emit 'start', @name
    new ObservableProcess(@service-config.setup,
                          cwd: @config.root,
                          verbose: yes,
                          console: {log: @_log, error: @_log},
                          on-exit: ~> @emit('finished', @name) ; done!)


  _log: (text) ~>
    @emit 'output', {@name, text}




module.exports = ServiceInstaller
