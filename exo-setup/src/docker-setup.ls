require! {
  'fs-extra' : fs
  '../../exosphere-shared' : {global-exosphere-directory}
  'handlebars' : Handlebars
  'js-yaml' : yaml
  'path'
  'prelude-ls' : {Obj, map}
  'os'
}


# Renders docker-compose.yml file with service configuration
class DockerSetup

  ({@app-name, @role, @logger, @service-location, @docker-image}) ->
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
      EXOCOM_HOST: 'exocom'
      EXOCOM_PORT: '$EXOCOM_PORT'
    for dependency of @service-config.dependencies
      env-vars[dependency.to-upper-case!] = dependency
    env-vars


  # compiles list of names of dependencies a service relies on
  _get-service-dependencies: ->
    dependencies = ['exocom']
    for dependency, dependency-config of @service-config.dependencies
      dependencies.push (dependency + dependency-config.dev.version)
    dependencies


  # builds the Docker config for a service dependency
  _get-service-dependency-docker-config: (dependency-name, dependency-config) ->
    if dependency-config.volumes
      data-path = global-exosphere-directory dependency-name 
      fs.ensure-dir-sync data-path
      rendered-volumes =  map ((volume) -> Handlebars.compile(volume)({"EXO_DATA_PATH": data-path})), dependency-config.volumes

    Obj.compact do
      image: "#{dependency-config.image}:#{dependency-config.version}"
      container_name: dependency-name + dependency-config.version
      ports: dependency-config.ports 
      volumes: rendered-volumes or undefined
    

  _get-external-service-docker-config: ->
    | !@docker-image => throw new Error red "No location or docker-image specified"
    docker-config = {}
    docker-config[@role] =
      image: @docker-image
      container_name: @role
      depends_on: ['exocom']
    docker-config


module.exports = DockerSetup
