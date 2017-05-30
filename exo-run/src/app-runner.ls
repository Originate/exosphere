require! {
  'chalk': {red}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerCompose}
  'path'
  './service-watcher' : ServiceWatcher
}


# Runs the overall application
class AppRunner extends EventEmitter

  ({@app-config, @logger}) ->
    @env =
      EXOCOM_PORT: process.env.EXOCOM_PORT or 80
    @docker-config-location = path.join process.cwd!, 'tmp'


  start: ->
    @watch-services!
    DockerHelper.run-all-images {@env, cwd: @docker-config-location, @write}, (exit-code) ~>
      | exit-code => return @shutdown error-message: 'Failed to run images'


  watch-services: ->
    @services = []
    for protection-level of @app-config.services
      for role, service-data of @app-config.services[protection-level]
        if service-data.location
          new ServiceWatcher {role, service-location: path.join(process.cwd!, service-data.location), @env, @logger}
            ..watch!
            ..on 'error', (message) ~> @shutdown error-message: message


  shutdown: ({close-message, error-message}) ~>
    switch
      | error-message  =>  console.log red error-message; exit-code = 1
      | otherwise      =>  console.log "\n\n #{close-message}"; exit-code = 0
    DockerHelper.kill-all-containers {cwd: @docker-config-location, @write}, -> process.exit exit-code


  write: (text) ~>
    @logger.log {role: 'exo-run', text, trim: yes}


module.exports = AppRunner
