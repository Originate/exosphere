require! {
  './app-linter' : AppLinter
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'fs'
}

module.exports = ->

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  console.log "Running linter for #{cyan app-config.name}\n"
  logger = new Logger Object.keys(app-config.services)
  app-linter = new AppLinter {app-config, logger}
    ..start!
