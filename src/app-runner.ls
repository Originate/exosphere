require! {
  'chalk' : {red}
  'events' : {EventEmitter}
  'exocomm-dev' : ExoComm
}


class AppRunner extends EventEmitter

  ->
    @services = {}


  start-exocomm: (port, done) ->
    exocomm = new ExoComm
      ..on 'error', (err) -> console.log red err
      ..on 'listening', (port) ~>
        console.log 2222222222
        @emit 'exocomm-online', port
      ..listen port



module.exports = AppRunner
