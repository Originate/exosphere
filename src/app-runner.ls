require! {
  'async'
  'chalk': {red}
  'child_process'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  'exocom-dev' : ExoCom
  'nitroglycerin' : N
  'port-reservation'
  'path'
  'require-yaml'
  './service-runner' : ServiceRunner
  'wait' : {wait-until}
}


# Runs the overall application
class AppRunner extends EventEmitter

  (@app-config) ->


  start-exocom: (done) ->
    port-reservation
      ..get-port N (@exocom-port) ~>
        service-messages = @compile-service-messages! |> JSON.stringify |> (.replace /"/g, '')
        @docker-config =
          author: 'originate'
          image: 'exocom'
          start-command: 'bin/exocom'
          env:
            SERVICE_MESSAGES: service-messages
            PORT: @exocom-port
            SERVICE_NAME: 'exocom'
          publish:
            EXOCOM_PORT: "#{@exocom-port}:#{@exocom-port}"
        @exocom = new DockerRunner 'exocom', @docker-config
          ..start-service!
          ..on 'output', (data) ~> @emit 'output', data
          ..on 'online', ~>
            @emit 'exocom-online', @exocom-port
          ..on 'error', (message) ~> @shutdown error-message: message


  start-services: ->
    wait-until (~> @exocom-port), 1, ~>
      names = Object.keys @app-config.services
      @runners = {}
      for name in names
        service-dir = path.join process.cwd!, @app-config.services[name].location
        @runners[name] = new ServiceRunner name, root: service-dir, EXOCOM_PORT: @exocom-port
          ..on 'online', (name) ~> @emit 'service-online', name
          ..on 'output', (data) ~> @emit 'output', data
          ..on 'error', @shutdown
      async.parallel [runner.start for _, runner of @runners], (err) ~>
        @emit 'all-services-online'


  shutdown: ({close-message, error-message}) ~>
    child_process.exec 'docker rm -f exocom'
    for service in Object.keys @app-config.services
      child_process.exec "docker rm -f #{service}"
    switch
      | error-message  =>  console.log red error-message; process.exit 1
      | otherwise      =>  console.log "\n\n #{close-message}"; process.exit!


  compile-service-messages: ->
    config = for service-name, service-data of @app-config.services
      service-config = require path.join(process.cwd!, service-data.location, 'service.yml')
      {
        name: service-name
        receives: service-config.messages.receives
        sends: service-config.messages.sends
        namespace: service-config.messages.namespace
      }



module.exports = AppRunner
