require! {
  'async'
  'events' : {EventEmitter}
  'path'
  './service-setup' : ServiceSetup
}


class AppSetup extends EventEmitter

  (@app-config) ->


  start-setup: ->
    setups = for own service-name of @app-config.services
      new ServiceSetup service-name, root: path.join(process.cwd!, @app-config.services[service-name].location)
        ..on 'start', (name) ~> @emit 'start', name
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'finished', (name, exit-code) ~> @emit 'finished', name, exit-code
        ..on 'error', (name, exit-code) ~>   @emit 'error', name, exit-code
    async.map setups, (-> &0.start &1), (err) ~>
      @emit 'setup-complete'



module.exports = AppSetup
