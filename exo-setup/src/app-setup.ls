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
    for protection-level of @app-config.services
      for service-role, service-data of @app-config.services[protection-level]
        @services.push do
            role: service-role
            location: service-data.location
    setups = for service in @services
      new ServiceSetup role: service.role, logger: @logger, config: root: path.join(process.cwd!, service.location)
        ..on 'output', (data) ~> @emit 'output', data

    docker-setups = for service in @services
      new DockerSetup role: service.role, logger: @logger, config: root: path.join(process.cwd!, service.location)
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
