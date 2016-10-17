require! {
  'async'
  './docker-setup' : DockerSetup
  'events' : {EventEmitter}
  'path'
  './service-setup' : ServiceSetup
}


class AppSetup extends EventEmitter

  ({@app-config, @logger}) ->


  start-setup: ->
    setups = for service-name of @app-config.services
      new ServiceSetup name: service-name, logger: @logger, config: root: path.join(process.cwd!, @app-config.services[service-name].location)
        ..on 'output', (data) ~> @emit 'output', data

    docker-setups = for service-name of @app-config.services
      new DockerSetup name: service-name, logger: @logger, config: root: path.join(process.cwd!, @app-config.services[service-name].location)
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



module.exports = AppSetup
