require! {
  'async'
  'chalk': {red}
  'child_process'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  '../../exosphere-shared' : {compile-service-messages, DockerHelper}
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
        service-messages = compile-service-messages @app-config |> JSON.stringify |> (.replace /"/g, '')
        @docker-config =
          author: 'originate'
          image: 'exocom'
          app-name: @app-config.name
          start-command: 'bin/exocom'
          env:
            SERVICE_ROUTES: service-messages
            PORT: @exocom-port
            ROLE: 'exocom'
          publish:
            EXOCOM_PORT: "#{@exocom-port}:#{@exocom-port}"
        @exocom = new DockerRunner {name: 'exocom', @docker-config, @logger}
          ..start-service!
          ..on 'error', (message) ~> @shutdown error-message: message


  start-services: ->
    wait-until (~> @exocom-port), 1, ~>
      @services = []
      for service-type of @app-config.services
        for service-name, service-data of @app-config.services[service-type]
          @services.push do
            {
              name: service-name
              location: service-data.location
              image: service-data.docker_image
            }
      @runners = {}
      for service in @services
        @runners[service.name] = new ServiceRunner {service.name, config: {root: path.join(process.cwd!, service.location), EXOCOM_PORT: @exocom-port, image: service.image, app-name: @app-config.name}, @logger}
          ..on 'error', @shutdown
      async.parallel [runner.start for _, runner of @runners], (err) ~>
        @logger.log name: 'exo-run', text: 'all services online'


  shutdown: ({close-message, error-message}) ~>
    switch
      | error-message  =>  console.log red error-message; exit-code = 1
      | otherwise      =>  console.log "\n\n #{close-message}"; exit-code = 0
    DockerHelper.remove-container \exocom
    for service in @services
      DockerHelper.remove-container service.name
      @runners[service.name].shutdown-dependencies!
    process.exit exit-code


module.exports = AppRunner
