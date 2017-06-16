require! {
  'fs-extra' : fs
  '../../exosphere-shared' : {ApplicationDependency, global-exosphere-directory, DockerHelper}
  'handlebars' : Handlebars
  'js-yaml' : yaml
  'path'
  'prelude-ls' : {Obj, map}
  'os'
}


# Renders docker-compose.yml file with service configuration
class DockerSetup

  ({@app-config, @role, @logger, @service-location, @docker-image}) ->
    @service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, @service-location, 'service.yml'), 'utf8') if @service-location


  get-service-docker-config: (done) ~>
    | @service-config => done null, @_get-service-docker-config!
    | otherwise       => @_get-external-service-docker-config done


  # builds the Docker config for a service and its dependencies
  _get-service-docker-config: ->
    docker-config = {}
    docker-config[@role] = Obj.compact do
      build: path.join '..', @service-location
      container_name: @role
      command: @service-config.startup.command
      ports: @service-config.docker?.ports or undefined
      links: @_get-docker-links!
      environment: @_get-docker-env-vars!
      depends_on: @_get-service-dependencies!
    for dependency, dependency-config of @service-config.dependencies
      docker-config[dependency + dependency-config.dev.version] = @_get-service-dependency-docker-config dependency, dependency-config.dev
    docker-config


  # compiles list of container links from service to dependency
  # returns undefined if length is 0 so it can be ignored with Obj.compact
  _get-docker-links: ->
    links = []
    for dependency, dependency-config of @service-config.dependencies
      links.push "#{dependency + dependency-config.dev.version}:#{dependency}"
    if links.length then links else undefined


  # compiles hash of environment variables to be set in a service container
  _get-docker-env-vars: ->
    env-vars =
      ROLE: @role
    for dependency-config in @app-config.dependencies
      dependency = ApplicationDependency.build dependency-config
      env-vars = {...env-vars, ...dependency.get-service-env-variables!}
    for dependency of @service-config.dependencies
      env-vars[dependency.to-upper-case!] = dependency
    env-vars


  # compiles list of names of dependencies a service relies on
  _get-service-dependencies: ->
    dependencies = @_get-app-dependencies!
    for dependency, dependency-config of @service-config.dependencies
      dependencies.push "#{dependency}#{dependency-config.dev.version}"
    dependencies


  # builds the Docker config for a service dependency
  _get-service-dependency-docker-config: (dependency-name, dependency-config) ->
    Obj.compact do
      image: "#{dependency-config.image}:#{dependency-config.version}"
      container_name: dependency-name + dependency-config.version
      ports: dependency-config.ports
      volumes: @_get-rendered-volumes dependency-config.volumes, dependency-name


  _get-external-service-docker-config: (done) ~>
    | !@docker-image => done new Error red "No location or docker-image specified"
    DockerHelper.cat-file image: @docker-image, file-name: 'service.yml', (err, external-service-config) ~>
      | err => done new Error red "Could not find the configuration for the docker-image"
      @service-config = yaml.safe-load external-service-config
      docker-config = {}
      docker-config[@role] = Obj.compact do
        image: @docker-image
        container_name: @role
        ports: @service-config.docker.ports
        environment: {...@service-config.docker.environment, ...@_get-docker-env-vars!}
        volumes: @_get-rendered-volumes @service-config.docker.volumes, @role
        depends_on: @_get-service-dependencies!
      done null, docker-config


  _get-rendered-volumes: (volumes, role)->
    if volumes
      data-path = global-exosphere-directory @app-config.name, role
      fs.ensure-dir-sync data-path
      map ((volume) -> Handlebars.compile(volume)({"EXO_DATA_PATH": data-path})), volumes


  _get-app-dependencies: ->
    map ((dependency-config) -> "#{dependency-config.type}#{dependency-config.version}"), @app-config.dependencies


module.exports = DockerSetup
