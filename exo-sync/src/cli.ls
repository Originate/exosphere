require! {
  'chalk' : {cyan}
  'fs'
  'js-yaml': yaml
  './app-syncer' : AppSyncer
  '../../exosphere-shared' : {Logger}
}

module.exports = ->

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf-8')
  logger = new Logger Object.keys(app-config.services ? {})
    ..log role: 'exo-sync', text: "Syncing application #{cyan app-config.name}"
  syncer = new AppSyncer {app-config, logger}
    ..start-syncing!
