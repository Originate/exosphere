require! {
  './app-setup' : AppSetup
  'chalk' : {cyan, green}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'fs'
  'prelude-ls' : {flatten}
}

module.exports = ->

  if process.argv[2] is \help
    return help!

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  console.log "Setting up #{green app-config.name} #{cyan app-config.version}\n"
  services = []
  silenced-services = []
  for type of app-config.services
    for service of app-config.services[type]
      services.push service
      if app-config.services[type][service].silent
        silenced-services.push service
  silenced-dependencies = [dependency.type if dependency.silent for dependency in app-config.dependencies]
  logger = new Logger services, silenced-services ++ silenced-dependencies
  app-setup = new AppSetup app-config: app-config, logger: logger
    ..start-setup!


function help
  help-message =
    """
    Usage: #{cyan "exo setup"}

    Sets up an Exosphere application so that it is ready to be run.
    The setup process includes:
      - locally installing dependencies of all application services
      - generating docker images for the application services
      - downloading the Exocom Docker image

    This command must be run in the root directory of the application.
    """
  console.log help-message
