require! {
  'chalk' : {red}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
  './service-runner' : ServiceRunner
}


class AppRunner extends EventEmitter

  ->
    @services = {}


  start-exocomm: (port, done) ->
    exocomm = new ExoComm
      ..on 'error', (err) -> console.log red err
      ..on 'listening', (port) ~>
        @emit 'exocomm-online', port
      ..listen port


  start-service: (name, config) ->
    new ServiceRunner name, config
      ..run!




module.exports = AppRunner
