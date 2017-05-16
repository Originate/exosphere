require! {
  'async'
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
      @logger.log role: 'exo-test', text: "#{@role} has no tests, skipping"
      return done?!

    @_start-dependencies (err) ~>
      | err  =>  @logger.log role: 'exo-test', text: "Error: #{err}"; return done?(err, 1)
      new ObservableProcess(call-args(@_create-command @service-config.tests)
                            cwd: @config.root,
                            env: @config
                            stdout: {@write}
                            stderr: {@write})
        ..on 'ended', (exit-code) ~>
          if exit-code > 0
            @logger.log role: 'exo-test', text: "#{@role} is broken"
          else
            @logger.log role: 'exo-test', text: "#{@role} works"
          @remove-dependencies!
          done?(null, exit-code)


  remove-dependencies: ~>
    for dep of @service-config.dependencies
      DockerHelper.remove-container "test-#{dep}"


  _start-dependencies: (done) ~>
    dependencies = []
    for dependency-name, dependency-config of @service-config.dependencies
      container-name = "test-#{dependency-name}"
      if dependency-config?.version?
        dependency-name += ":#{dependency-config.version}"
      if dependency-config?.docker_flags?
        online-text = that.online_text
        port = that.port
      dependencies.push {container-name, dependency-name, online-text, port}
    async.each-series dependencies, DockerHelper.ensure-container-is-running, (err) ~>
      | err  => done err
      done!


  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @config.root, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  write: (text) ~>
    @logger.log {@role, text}



module.exports = ServiceTester
