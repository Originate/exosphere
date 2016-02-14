require! {
  'async'
  'events' : {EventEmitter}
  'path'
  './service-installer' : ServiceInstaller
}


class AppInstaller extends EventEmitter

  (@app-config) ->


  start-installation: ->
    installers = for service-name in Object.keys @app-config.services
      new ServiceInstaller service-name, root: path.join(process.cwd!, @app-config.services[service-name].location)
        ..on 'start', (name) ~> @emit 'start', name
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'finished', (name) ~> @emit 'finished', name
    async.parallel [installer.start for installer in installers], (err) ~>
      @emit 'installation-complete'



module.exports = AppInstaller
