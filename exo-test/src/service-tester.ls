require! {
  'events' : {EventEmitter}
  '../../exosphere-shared' : {call-args, DockerHelper}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
  'wait' : {wait}
}


class ServiceTester extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    unless @service-config.tests
      @emit 'service-tests-skipped', @name
      return done?!

    @_start-dependencies (err) ~>
      | err => @emit 'error', err ; return
      new ObservableProcess(call-args(@_create-command @service-config.tests)
                            cwd: @config.root,
                            env: @config
                            stdout: {@write}
                            stderr: {@write})
        ..on 'ended', (exit-code) ~>
          if exit-code > 0
            @emit 'service-tests-failed', @name
          else
            @emit 'service-tests-passed', @name
          @remove-dependencies!
          done?(null, exit-code)


  remove-dependencies: ~>
    if @service-config.dependencies
      for dep in @service-config.dependencies
        DockerHelper.remove-container "test-#{dep}"


  _start-dependencies: (done) ~>
    if @service-config.dependencies
      for dep in @service-config.dependencies
        DockerHelper.ensure-container-is-running "test-#{dep}", dep
      wait 500, done
    else
      done!


  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @config.root, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  write: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceTester
