require! {
  'async'
  'chalk': {red}
  'child_process'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  '../../exosphere-shared' : {compile-service-routes, DockerHelper}
  'nitroglycerin' : N
  'port-reservation'
  'path'
  './service-runner' : ServiceRunner
  'wait' : {wait-until}
}


# Runs the overall application
class AppRunner extends EventEmitter

  ({@app-config, @logger}) ->


  start-exocom: (done) ->
    port-reservation
      ..get-port N (@exocom-port) ~>
        service-routes = compile-service-routes @app-config |> JSON.stringify |> (.replace /"/g, '')
        @docker-config =
          author: 'originate'
          image: @app-config.bus.type
          version: @app-config.bus.version
          app-name: @app-config.name
          start-command: 'bin/exocom'
          start-text: 'ExoCom WebSocket listener online at port'
          env:
            SERVICE_ROUTES: service-routes
            PORT: @exocom-port
            ROLE: 'exocom'
          publish:
            EXOCOM_PORT: "#{@exocom-port}:#{@exocom-port}"
        @exocom = new DockerRunner {role: 'exocom', @docker-config, @logger}
          ..start-service!
          ..on 'error', (message) ~> @shutdown error-message: message


  start-services: ->
    wait-until (~> @exocom-port), 1, ~>
      @services = []
      for protection-level of @app-config.services
        for role, service-data of @app-config.services[protection-level]
          @services.push do
            {
              role: role
              location: service-data.location
              image: service-data.docker_image
            }
      @runners = {}
      for service in @services
        @runners[service.role] = new ServiceRunner {service.role, config: {root: path.join(process.cwd!, service.location), EXOCOM_PORT: @exocom-port, image: service.image, app-name: @app-config.name}, @logger}
          ..on 'error', @shutdown
      async.parallel [runner.start for _, runner of @runners], (err) ~>
        @logger.log role: 'exo-run', text: 'all services online'


  shutdown: ({close-message, error-message}) ~>
    switch
      | error-message  =>  console.log red error-message; exit-code = 1
      | otherwise      =>  console.log "\n\n #{close-message}"; exit-code = 0
    DockerHelper.remove-container \exocom
    for service in @services
      DockerHelper.remove-container service.role
      @runners[service.role].shutdown-dependencies!
    process.exit exit-code


module.exports = AppRunner
