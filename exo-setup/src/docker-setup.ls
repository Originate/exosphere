require! {
  'handlebars' : Handlebars
  'fs-extra' : fs
  'prelude-ls' : {Obj, map}
  'js-yaml' : yaml
  'path'
  'os'
}


# Renders docker-compose.yml file with service configuration
class DockerSetup

  ({@app-name, @role, @logger, @service-location, @docker-image}) ->
    @service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, @service-location, 'service.yml'), 'utf8') if @service-location


  get-service-docker-config: ~> 
    | @service-config => @_setup-service!
    | otherwise       => @_setup-external-service!


  _setup-service: ->
    docker-config = {}
    docker-config[@role] =
      Obj.compact do
        build: @service-location 
        command: @service-config.startup.command
        ports: @service-config.docker?.ports or null
        links: @_get-service-links!
        environment: @_get-environment-vars! 
        depends_on: @_get-service-dependencies!
    for dependency, dependency-config of @service-config.dependencies
      docker-config[dependency + dependency-config.dev.version] = @_setup-service-dependencies dependency, dependency-config.dev
    docker-config


  _get-service-links: ->
    links = []
    for dependency, dependency-config of @service-config.dependencies
      links.push "#{dependency + dependency-config.dev.version}:#{dependency}"
    if links.length then links else null


  _get-environment-vars: ->
    env-vars =
      ROLE: @role
      EXOCOM_HOST: 'exocom'
      EXOCOM_PORT: '$EXOCOM_PORT'
    for dependency of @service-config.dependencies
      env-vars[dependency.to-upper-case!] = dependency
    env-vars


  _get-service-dependencies: ->
    dependencies = ['exocom']
    for dependency, dependency-config of @service-config.dependencies
      dependencies.push (dependency + dependency-config.dev.version)
    dependencies


  _setup-service-dependencies: (dependency-name, dependency-config) ->
    if dependency-config.volumes
      data-path = path.join os.homedir!, '.exosphere', @app-name, dependency-name, 'data'
      fs.ensure-dir-sync data-path
      rendered-volumes =  map ((volume) -> Handlebars.compile(volume)({"EXO_DATA_PATH": data-path})), dependency-config.volumes

    dependency =
      image: "#{dependency-config.image}:#{dependency-config.version}"
      ports: dependency-config.ports 
      volumes: rendered-volumes or null
    Obj.compact dependency
    
  _setup-external-service: ->
    throw new Error red "No location or docker-image specified" unless @docker-image
    docker-config = {}
    docker-config[@role] =
      image: @docker-image
      depends_on: ['exocom']
    docker-config
    
module.exports = DockerSetup
