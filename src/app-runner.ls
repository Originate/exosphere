require! {
  'async'
  'chalk' : {red}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
  './next-port'
  './service-runner' : ServiceRunner
}


# Runs the overall application
class AppRunner extends EventEmitter

  (@app-config) ->


  start-exocomm: (done) ->
    next-port (@exocomm-port) ~>
      @exocomm = new ExoComm
        ..on 'error', (err) -> console.log red err
        ..on 'listening', (port) ~> @emit 'exocomm-online', port
        ..listen @exocomm-port


  start-services: ->
    names = Object.keys @app-config.services
    @runners = {}
    for name in names
      @runners[name] = new ServiceRunner name, 'exocomm-port': @exocomm-port
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
        port: runner.config['exorelay-port']
        sends: runner.service-config.messages.sends
        receives: runner.service-config.messages.receives
      }
    @exocomm.set-services config
    @emit 'routing-done'



module.exports = AppRunner
