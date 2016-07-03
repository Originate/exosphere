require! {
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


class ServiceTester extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    new ObservableProcess(@_create-command(@service-config.tests)
                          cwd: @config.root,
                          env: @config
                          verbose: yes,
                          console: log: @_log, error: @_log)
      ..on 'ended', (exit-code) ~>
        @emit 'done', @name
        done null, exit-code


  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @config.root, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  _log: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceTester
