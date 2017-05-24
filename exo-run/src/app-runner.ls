require! {
  'chalk': {red}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper}
  'path'
  './service-watcher' : ServiceWatcher
}


# Runs the overall application
class AppRunner extends EventEmitter

  ({@app-config, @logger}) ->
    @env =
      EXOCOM_PORT: 80 or process.env.EXOCOM_PORT


  start: ->
    @watch-services!
    DockerHelper.run-all-images {@env, @write}, (exit-code, killed) ~>
      | exit-code => return @shutdown error-message: 'Failed to run images'
      @logger.log role: 'exo-run', text: 'all services online'


  watch-services: ->
    @services = []
    for protection-level of @app-config.services
      for role, service-data of @app-config.services[protection-level]
        if service-data.location
          new ServiceWatcher {role, service-location: path.join(process.cwd!, service-data.location), @env, @logger}
            ..watch!


  shutdown: ({close-message, error-message}) ~>
    switch
      | error-message  =>  console.log red error-message; exit-code = 1
      | otherwise      =>  console.log "\n\n #{close-message}"; exit-code = 0
    DockerHelper.kill-all-containers {@write}, -> process.exit exit-code


  write: (text) ~>
    @logger.log {role: 'exo-run', text, trim: yes}


module.exports = AppRunner
