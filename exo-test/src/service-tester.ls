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

  ({@role, @config, @logger}) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    unless @service-config.tests
      @logger.log name: 'exo-test', text: "#{@role} has no tests, skipping"
      return done?!

    @_start-dependencies (err) ~>
      | err => @logger.log name: 'exo-test', text: "Error: #{err}"; return
      new ObservableProcess(call-args(@_create-command @service-config.tests)
                            cwd: @config.root,
                            env: @config
                            stdout: {@write}
                            stderr: {@write})
        ..on 'ended', (exit-code) ~>
          if exit-code > 0
            @logger.log name: 'exo-test', text: "#{@role} is broken"
          else
            @logger.log name: 'exo-test', text: "#{@role} works"
          @remove-dependencies!
          done?(null, exit-code)


  remove-dependencies: ~>
    for dep of @service-config.dependencies
      DockerHelper.remove-container "test-#{dep}"


  _start-dependencies: (done) ~>
    for dep of @service-config.dependencies
      DockerHelper.ensure-container-is-running "test-#{dep}", dep
    wait 500, done


  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @config.root, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  write: (text) ~>
    @logger.log {name: @role, text}



module.exports = ServiceTester
