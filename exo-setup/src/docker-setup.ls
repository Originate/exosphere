require! {
  'fs-extra' : fs
  '../../exosphere-shared' : {ApplicationDependency, global-exosphere-directory}
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


  get-service-docker-config: ~>
    | @service-config => @_get-service-docker-config!
    | otherwise       => @_get-external-service-docker-config!


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
    if @service-config.dependencies
      for dependency in @service-config.dependencies
        console.log 'dependency' dependency
        docker-config[dependency.name + dependency.version] = @_get-service-dependency-docker-config dependency.name, dependency.version, dependency.config
    docker-config


  # compiles list of container links from service to dependency
  # returns undefined if length is 0 so it can be ignored with Obj.compact
  _get-docker-links: ->
    links = []
    # console.log @service-config.dependencies
    if @service-config.dependencies
      for dependency in @service-config.dependencies
        links.push "#{dependency.name + dependency.version}:#{dependency.name}"
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
    if @service-config.dependencies
      for dependency in @service-config.dependencies
        dependencies.push "#{dependency.name}#{dependency.version}"
      dependencies


  # builds the Docker config for a service dependency
  _get-service-dependency-docker-config: (dependency-name, dependency-version, dependency-config) ->
    console.log 'dependency-name' dependency-name
    console.log 'dependency-config' dependency-config
    if dependency-config.volumes
      data-path = global-exosphere-directory @app-config.name, dependency-name
      fs.ensure-dir-sync data-path
      rendered-volumes =  map ((volume) -> Handlebars.compile(volume)({"EXO_DATA_PATH": data-path})), dependency-config.volumes

    Obj.compact do
      image: "#{dependency-name}:#{dependency-version}"
      container_name: dependency-name + dependency-version
      ports: dependency-config.ports
      volumes: rendered-volumes or undefined


  _get-external-service-docker-config: ->
    | !@docker-image => throw new Error red "No location or docker-image specified"
    docker-config = {}
    docker-config[@role] =
      image: @docker-image
      container_name: @role
      depends_on: @_get-app-dependencies!
    docker-config


  _get-app-dependencies: ->
    map ((dependency-config) -> "#{dependency-config.name}#{dependency-config.version}"), @app-config.dependencies


module.exports = DockerSetup
