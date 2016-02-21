require! {
  'async'
  'chalk' : {red}
  'rails-delegate' : {delegate-event}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
  'nitroglycerin' : N
  'port-reservation'
  'path'
  './service-runner' : ServiceRunner
  'wait' : {wait-until}
}


# Runs the overall application
class AppRunner extends EventEmitter

  (@app-config) ->


  start-exocomm: (done) ->
    port-reservation.get-port N (@exocomm-port) ~>
      @exocomm = new ExoComm
        ..on 'listening', (port) ~> @emit 'exocomm-online', port
        ..listen @exocomm-port
      delegate-event 'error', 'routing-setup', 'message', from: @exocomm, to: @


  start-services: ->
    wait-until (~> @exocomm-port), 1, ~>
      names = Object.keys @app-config.services
      @runners = {}
      for name in names
        service-dir = path.join process.cwd!, @app-config.services[name].location
        @runners[name] = new ServiceRunner name, root: service-dir, EXOCOMM_PORT: @exocomm-port
          ..on 'online', (name) ~> @emit 'service-online', name
          ..on 'output', (data) ~> @emit 'output', data
      async.parallel [runner.start for _, runner of @runners], (err) ~>
        @emit 'all-services-online'


  # Sends which service listens on what port to ExoComm
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
    @exocomm.set-services config
    @emit 'routing-done'



module.exports = AppRunner
