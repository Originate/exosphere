require! {
  'async'
  'asynchronizer' : Asynchronizer
  'chalk': {red}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {ApplicationDependency, DockerCompose, DockerHelper}
  'fs'
  'path'
  './service-restarter' : ServiceRestarter
  'js-yaml' : yaml
}


# Runs the overall application
class AppRunner extends EventEmitter

  ({@app-config, @logger}) ->
    @env = {}
    for dependency-config in @app-config.dependencies
      dependency = ApplicationDependency.build dependency-config
      @env = {...@env, ...dependency.get-env-variables!}
    @docker-config-location = path.join process.cwd!, 'tmp'


  start: ->
    @watch-services!
    @process = DockerCompose.run-all-images {@env, cwd: @docker-config-location, @write}, (exit-code) ~>
      | exit-code => return @shutdown error-message: 'Failed to run images'

    @_compile-online-text (err) ~>
      | err => throw err
      asynchronizer = new Asynchronizer Object.keys(@online-texts)
      for role, online-text of @online-texts
        let role, online-text
          @process.wait (new RegExp(role + ".*" + online-text)), ~>
            @logger.log {role, text: "'#{role}' is running"}
            asynchronizer.check role
      asynchronizer.then ~>
        @write 'all services online'


  watch-services: ->
    @services = []
    for protection-level of @app-config.services
      for role, service-data of @app-config.services[protection-level]
        if service-data.location
          new ServiceRestarter {role, service-location: path.join(process.cwd!, service-data.location), @env, @logger}
            ..watch!
            ..on 'error', (message) ~> @shutdown error-message: message


  shutdown: ({close-message, error-message}) ~>
    switch
      | error-message  =>  console.log red error-message; exit-code = 1
      | otherwise      =>  console.log "\n\n #{close-message}"; exit-code = 0
    DockerCompose.kill-all-containers {cwd: @docker-config-location, @write}, -> process.exit exit-code


  _compile-online-text: (done) ~>
    @online-texts = {}
    for app-dependency in @app-config.dependencies
      dependency = ApplicationDependency.build app-dependency
      @online-texts[app-dependency.type] = dependency.get-online-text!
    services = []
    for protection-level of @app-config.services
      for role, service-data of @app-config.services[protection-level]
        services.push {role: role, service-data: service-data}
    async.map-series services, @_get-online-text, (err) ~>
      | err => done err
      done!


  _get-online-text: ({role, service-data}, done) ~>
    | service-data.location
      service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, service-data.location, 'service.yml'))
      @online-texts[role] = service-config.startup['online-text']
      done!
    | service-data['docker-image'] =>
      DockerHelper.cat-file image: service-data['docker-image'], file-name: 'service.yml', (err, external-service-config) ~>
        | err => done new Error red "Could not find the configuration for the docker-image"
        service-config = yaml.safe-load external-service-config
        @online-texts[role] = service-config.startup['online-text']
        done!
    | otherwise => done new Error red "No location or docker-image specified"


  write: (text) ~>
    @logger.log {role: 'exo-run', text, trim: yes}


module.exports = AppRunner
