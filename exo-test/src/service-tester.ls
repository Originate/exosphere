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

  ({@role, @service-location, @logger}) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@service-location, 'service.yml'), 'utf8')


  start: (done) ~>
    unless @service-config.tests
      @logger.log role: 'exo-test', text: "#{@role} has no tests, skipping"
      return done?!

    @_start-dependencies (err) ~>
      | err  =>  @logger.log role: 'exo-test', text: "Error: #{err}"; return done?(err, 1)
      new ObservableProcess(call-args(@_create-command @service-config.tests)
                            cwd: @service-location,
                            stdout: {@write}
                            stderr: {@write})
        ..on 'ended', (exit-code) ~>
          if exit-code > 0
            @logger.log role: 'exo-test', text: "#{@role} is broken"
          else
            @logger.log role: 'exo-test', text: "#{@role} works"
          @_remove-dependencies (err) ->
            done?(err, exit-code)


  _remove-dependencies: (done) ~>
    async.each-series @dependencies, DockerHelper.remove-container, (err) ~>
      | err  => done err
      done!


  _start-dependencies: (done) ~>
    @dependencies = []
    for dependency-name, dependency-config of @service-config.dependencies
      if dependency-config.dev
        @dependencies.push do
          Image: "#{that.image}:#{that.version}"
          name: "#{@role}-test-#{dependency-name}"
          HostConfig: @_get-port-mapping that 
          online-text: that['online-text']
    async.each-series @dependencies, DockerHelper.start-container, (err) ~>
      | err  => done err
      done!


  # Converts port mappings from 'host_port:container:port' the object: 
  #
  # PortBindings: 
  #   "27017/tcp": [{HostPort: "27017"}]
  _get-port-mapping: (dependency-config) ->
    port-mappings = PortBindings: {}
    for port-mapping in dependency-config.ports
      host-port = port-mapping.split(':')[0]
      container-port = port-mapping.split(':')[1]
      port-mappings.PortBindings["#{container-port}/tcp"] = [{HostPort: host-port}]
    port-mappings
    

  _create-command: (command) ->
    if @_is-local-command command
      command = path.join @service-location, command
    command


  _is-local-command: (command) ->
    command.substr(0, 2) is './'


  write: (text) ~>
    @logger.log {@role, text}



module.exports = ServiceTester
