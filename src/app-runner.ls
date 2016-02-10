require! {
  'async'
  'chalk' : {red}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
  './service-runner' : ServiceRunner
}


# Runs the overall application
class AppRunner extends EventEmitter

  start-exocomm: (port, done) ->
    @exocomm-port = port
    exocomm = new ExoComm
      ..on 'error', (err) -> console.log red err
      ..on 'listening', (port) ~>
        @emit 'exocomm-online', port
      ..listen port


  start-services: (services) ->
    runners = for name, config of services
      config['exocomm-port'] = @exocomm-port
      new ServiceRunner(name, config)
    for runner in runners
      runner.on 'online', (name) ~> @emit 'service-online', name
      runner.on 'output', (data) ~> @emit 'output', data
    async.parallel [runner.start for runner in runners], (err) ~>
      @emit 'all-services-online'


module.exports = AppRunner
