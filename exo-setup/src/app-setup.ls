require! {
  'async'
  './docker-setup' : DockerSetup
  'events' : {EventEmitter}
  './exocom-setup' : ExocomSetup
  'path'
  './service-setup' : ServiceSetup
}


class AppSetup extends EventEmitter

  ({@app-config, @logger}) ->


  start-setup: ->
    @services = []
    for service-type of @app-config.services
      for service-name, service-data of @app-config.services[service-type]
        @services.push do
            name: service-name
            location: service-data.location
    setups = for service in @services
      new ServiceSetup name: service.name, logger: @logger, config: root: path.join(process.cwd!, service.location)
        ..on 'output', (data) ~> @emit 'output', data

    docker-setups = for service in @services
      new DockerSetup name: service.name, logger: @logger, config: root: path.join(process.cwd!, service.location)
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
        @logger.log name: 'exo-setup', text: 'setup complete'

    new ExocomSetup @logger
      ..start!



module.exports = AppSetup
