require! {
  'async'
  './docker-setup' : DockerSetup
  'events' : {EventEmitter}
  'path'
  './service-setup' : ServiceSetup
}


class AppSetup extends EventEmitter

  (@app-config) ->


  start-setup: ->
    setups = for service-name of @app-config.services
      new ServiceSetup service-name, root: path.join(process.cwd!, @app-config.services[service-name].location)
        ..on 'start', (name) ~> @emit 'start', name
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'finished', (name, exit-code) ~> @emit 'finished', name, exit-code
        ..on 'error', (name, exit-code) ~> @emit 'error', name, exit-code

    docker-setups = for service-name of @app-config.services
      new DockerSetup service-name, root: path.join(process.cwd!, @app-config.services[service-name].location)
        ..on 'docker-exists', (name) ~> @emit 'docker-exists', name
        ..on 'docker-finished', (name, exit-code) ~> @emit 'docker-finished', name, exit-code
        ..on 'docker-start', (name) ~> @emit 'docker-start', name
        ..on 'error', (name, exit-code) ~> @emit 'docker-error', name, exit-code
        ..on 'output', (data) ~> @emit 'output', data
    # Note: Windows does not provide atomic file operations,
    #       causing file system permission errors when multiple threads read and write to the same cache directory.
    #       We therefore run only one operation at a time on Windows.
    operation = if process.platform is 'win32'
      async.map-series
    else
      async.map
    operation setups, (-> &0.start &1), (err) ~>
      operation docker-setups, (-> &0.start &1), (err) ~>
        @emit 'setup-complete'



module.exports = AppSetup
