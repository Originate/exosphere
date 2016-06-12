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
        ..on 'finished', (name) ~> @emit 'finished', name
    async.map setups, (.start!), (err) ~>
      @emit 'setup-complete'



module.exports = AppSetup
