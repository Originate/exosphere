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

  start-exocomm: (done) ->
    next-port (@exocomm-port) ~>
      exocomm = new ExoComm
        ..on 'error', (err) -> console.log red err
        ..on 'listening', (port) ~> @emit 'exocomm-online', port
        ..listen @exocomm-port


  start-services: (services) ->
    names = [key for key, _ of services]
    runners = for name in names
      new ServiceRunner name, 'exocomm-port': @exocomm-port
    for runner in runners
      runner
        ..on 'online', (name) ~> @emit 'service-online', name
        ..on 'output', (data) ~> @emit 'output', data
    async.parallel [runner.start for runner in runners], (err) ~>
      @emit 'all-services-online'



module.exports = AppRunner
