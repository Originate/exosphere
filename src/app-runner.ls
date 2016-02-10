require! {
  'chalk' : {cyan, bg-cyan, magenta, bg-magenta, red}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
  './service-runner' : ServiceRunner
}


class AppRunner extends EventEmitter

  ->
    @services = {}
    @service-count = 0


  fg-colors: [cyan, magenta]
  bg-colors: [bg-cyan, bg-magenta]

  start-exocomm: (port, done) ->
    exocomm = new ExoComm
      ..on 'error', (err) -> console.log red err
      ..on 'listening', (port) ~>
        @emit 'exocomm-online', port
      ..listen port


  start-service: (name, config) ->
    new ServiceRunner name, config, @fg-colors[@service-count], @bg-colors[@service-count]
      ..run!
    @service-count++



module.exports = AppRunner
