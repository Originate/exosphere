require! {
  'lodash/assign' : assign
  'async'
  '../../exosphere-shared' : {DockerHelper, compile-service-routes}
  './docker-setup' : DockerSetup
  'events' : {EventEmitter}
  'fs-extra' : fs
  'path'
  './service-setup' : ServiceSetup
  'js-yaml' : yaml
}


class AppSetup extends EventEmitter

  ({@app-config, @logger}) ->
    @docker-compose-config =
      version: '3'
      services: {}
    @docker-compose-location = path.join process.cwd!, 'tmp', 'docker-compose.yml'


  start-setup: ->
    @services = []
    for protection-level of @app-config.services
      for service-role, service-data of @app-config.services[protection-level]
        @services.push do
            role: service-role
            location: service-data.location
            docker-image: service-data.docker_image

    @_setup-services ~>
      @_get-exocom-docker-config!
      @_get-service-docker-config!
      @_render-docker-compose!
      @_setup-docker-images (exit-code) ~>
        | exit-code => @write 'setup failed'; process.exit exit-code
        @write 'setup complete'


  _setup-services: (done) ->
    setups = for service in @services
      if service.location
        new ServiceSetup role: service.role, logger: @logger, config: root: path.join(process.cwd!, service.location)
          ..on 'output', (data) ~> @emit 'output', data
    async.map-series setups, (-> &0.start &1), (err) ~>
      | err  =>  throw new Error err
      done!


  _get-exocom-docker-config: ->
    exocom-docker-config =
      exocom:
        image: "originate/exocom:#{@app-config.bus.version}"
        command: 'bin/exocom'
        container_name: 'exocom'
        environment:
          ROLE: 'exocom'
          PORT: '$EXOCOM_PORT'
          SERVICE_ROUTES: compile-service-routes @app-config |> JSON.stringify |> (.replace /"/g, '')

    @docker-compose-config.services `assign` exocom-docker-config 


  _get-service-docker-config: ->
    docker-setups = for service in @services
      docker-setup = new DockerSetup do
        app-name: @app-config.name
        role: service.role
        logger: @logger
        service-location: service.location
        docker-image: service.docker-image
      @docker-compose-config.services `assign` docker-setup.get-service-docker-config!


  _render-docker-compose: ->
    fs.ensure-file-sync @docker-compose-location
    fs.write-file-sync @docker-compose-location, yaml.safe-dump(@docker-compose-config)


  _setup-docker-images: (done) ->
    DockerHelper.pull-all-images {@write, cwd: path.dirname @docker-compose-location}, (exit-code, killed) ~>
      | exit-code => @write 'Docker setup failed'; done exit-code
      | otherwise =>
        DockerHelper.build-all-images {@write, cwd: path.dirname @docker-compose-location}, (exit-code, killed) ~>
          | exit-code => @write 'Docker setup failed'; done exit-code
          | otherwise => @write 'Docker setup finished'; done!


  write: (text) ~>
    @logger.log {role: 'exo-setup', text, trim: yes}


module.exports = AppSetup
