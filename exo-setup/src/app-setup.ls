require! {
  'async'
  './docker-setup' : DockerSetup
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper, compile-service-routes}
  'fs-extra' : fs
  'lodash/assign' : assign
  'path'
  './service-setup' : ServiceSetup
  'js-yaml' : yaml
}


class AppSetup extends EventEmitter

  ({@app-config, @logger}) ->
    @docker-compose-config =
      version: '3'
      services: {}


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
      @_setup-docker-images (err) ~>
        | err => throw new Error err 
        @logger.log role: 'exo-setup', text: 'setup complete' 


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
    docker-compose-path = path.join(process.cwd!, 'docker-compose.yml') 
    fs.ensure-file-sync docker-compose-path 
    fs.write-file-sync docker-compose-path, yaml.safe-dump(@docker-compose-config)


  _setup-docker-images: (done) ->
    DockerHelper.pull-all-images {@write}, (exit-code, killed) ~>
      | exit-code => @logger.log role: 'exo-setup', text: 'Docker setup failed'; done exit-code
      | otherwise =>
        DockerHelper.build-all-images {@write}, (exit-code, killed) ~>
          | exit-code => @logger.log role: 'exo-setup', text: 'Docker setup failed'; done exit-code
          | otherwise => @logger.log role: 'exo-setup', text: 'Docker setup finished'; done!


  write: (text) ~>
    @logger.log {role: 'exo-setup', text, trim: yes}


module.exports = AppSetup
