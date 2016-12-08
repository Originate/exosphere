require! {
  './app-setup' : AppSetup
  'chalk' : {cyan, green}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'fs'
  'prelude-ls' : {flatten}
}

module.exports = ->

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  console.log "Setting up #{green app-config.name} #{cyan app-config.version}\n"
  logger = new Logger flatten [Object.keys(app-config.services[type]) for type of app-config.services]
  app-setup = new AppSetup app-config: app-config, logger: logger
    ..on 'output', (data) -> data.text = data.text.replace('\n', '') ; logger.log data
    ..start-setup!
