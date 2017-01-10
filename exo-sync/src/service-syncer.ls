require! {
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


class ServiceSyncer extends EventEmitter

  ({@role, @config, @logger}) ->


  start: (done) ~>
    | !@_is-external-service!  =>  return done null, 0

    new ObservableProcess("git pull",
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code is 0  =>  @write 'Sync ok'
        | otherwise       =>  @write 'Sync error'
        done null, exit-code


  _is-external-service: ->
    @config.root.split('/')[0] is '..'


  write: (text) ~>
    @logger.log {name: @role, text}



module.exports = ServiceSyncer
