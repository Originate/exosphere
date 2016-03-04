require! {
  'async'
  'chalk' : {red}
  'rails-delegate' : {delegate-event}
  'events' : {EventEmitter}
  'exocom-dev' : ExoCom
  'nitroglycerin' : N
  'port-reservation'
  'path'
  './service-runner' : ServiceRunner
  'wait' : {wait-until}
}


# Runs the overall application
class AppRunner extends EventEmitter

  (@app-config) ->


  start-exocom: (done) ->
    port-reservation.get-port N (@exocom-port) ~>
      @exocom = new ExoCom
        ..on 'listening', (port) ~> @emit 'exocom-online', port
        ..listen @exocom-port
      delegate-event 'error', 'routing-setup', 'message', from: @exocom, to: @


  start-services: ->
    wait-until (~> @exocom-port), 1, ~>
      names = Object.keys @app-config.services
      @runners = {}
      for name in names
        service-dir = path.join process.cwd!, @app-config.services[name].location
        @runners[name] = new ServiceRunner name, root: service-dir, EXOCOM_PORT: @exocom-port
          ..on 'online', (name) ~> @emit 'service-online', name
          ..on 'output', (data) ~> @emit 'output', data
      async.parallel [runner.start for _, runner of @runners], (err) ~>
        @emit 'all-services-online'


  # Returns the exorelay port for the service with the given name
  port-for: (service-name) ->
    @runners[service-name].config.EXORELAY_PORT


  # Sends which service listens on what port to ExoCom
  send-service-configuration: ->
    config = for service-name, service-data of @app-config.services
      runner = @runners[service-name]
      {
        name: service-name
        host: 'localhost'
        port: runner.config.EXORELAY_PORT
        sends: runner.service-config.messages.sends
        receives: runner.service-config.messages.receives
      }
    @exocom.set-services config
    @emit 'routing-done'



module.exports = AppRunner
