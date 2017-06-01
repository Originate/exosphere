require! {
  'asynchronizer' : Asynchronizer
  'chalk': {red}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerCompose}
  'fs'
  'path'
  './service-restarter' : ServiceRestarter
  'js-yaml' : yaml
}


# Runs the overall application
class AppRunner extends EventEmitter

  ({@app-config, @logger}) ->
    @env =
      EXOCOM_PORT: process.env.EXOCOM_PORT or 80
    @docker-config-location = path.join process.cwd!, 'tmp'


  start: ->
    @watch-services!
    @process = DockerCompose.run-all-images {@env, cwd: @docker-config-location, @write}, (exit-code) ~>
      | exit-code => return @shutdown error-message: 'Failed to run images'

    online-texts = @_compile-online-text!
    asynchronizer = new Asynchronizer Object.keys(online-texts)

    for role, online-text of online-texts 
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


  _compile-online-text: ->
    online-texts = {}
    for protection-level of @app-config.services
      for role, service-data of @app-config.services[protection-level]
        if service-data.location #TODO: compile online text for external services
          service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, service-data.location, 'service.yml'))
          online-texts[role] = service-config.startup['online-text']
    online-texts


  write: (text) ~>
    @logger.log {role: 'exo-run', text, trim: yes}


module.exports = AppRunner
