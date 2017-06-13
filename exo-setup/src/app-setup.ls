require! {
  'async'
  'events' : {EventEmitter}
  '../../exosphere-shared' : {ApplicationDependency, DockerCompose}
  './docker-setup' : DockerSetup
  'fs-extra' : fs
  'js-yaml' : yaml
  'lodash/assign' : assign
  'path'
  './service-setup' : ServiceSetup
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
            docker-image: service-data['docker-image']

    @_setup-services ~>
      @_get-dependencies-docker-config (err) ~>
        | err => @write 'setup failed'; process.exit 1
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


  _get-dependencies-docker-config: (done) ->
    console.log 'getting dependencies-docker-config'
    for dependency-config in @app-config.dependencies
      console.log dependency-config
      dependency = ApplicationDependency.build dependency-config
      console.log dependency
      dependency.get-docker-config @app-config, (err, docker-config) ~>
        | err => return done err
        @docker-compose-config.services `assign` docker-config
        done!


  _get-service-docker-config: ->
    console.log 'getting service-docker-config'
    docker-setups = for service in @services
      docker-setup = new DockerSetup {
        @app-config
        role: service.role
        @logger
        service-location: service.location
        docker-image: service.docker-image
      }
      @docker-compose-config.services `assign` docker-setup.get-service-docker-config!


  _render-docker-compose: ->
    console.log 'creating docker compose'
    console.log JSON.stringify @docker-compose-config
    fs.ensure-file-sync @docker-compose-location
    fs.write-file-sync @docker-compose-location, yaml.safe-dump(@docker-compose-config)


  _setup-docker-images: (done) ->
    DockerCompose.pull-all-images {@write, cwd: path.dirname @docker-compose-location}, (exit-code, killed) ~>
      | exit-code => @write 'Docker setup failed'; done exit-code
      | otherwise =>
        DockerCompose.build-all-images {@write, cwd: path.dirname @docker-compose-location}, (exit-code, killed) ~>
          | exit-code => @write 'Docker setup failed'; done exit-code
          | otherwise => @write 'Docker setup finished'; done!


  write: (text) ~>
    @logger.log {role: 'exo-setup', text, trim: yes}


module.exports = AppSetup
