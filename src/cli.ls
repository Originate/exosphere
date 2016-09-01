require! {
  'chalk' : {cyan}
  'fs'
  'js-yaml': yaml
  './app-syncer' : AppSyncer
  'exosphere-shared' : {Logger}
}


app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf-8')
logger = new Logger Object.keys(app-config.services ? {})
  ..log name: 'exo-sync', text: "Syncing application #{cyan app-config.name}"
syncer = new AppSyncer app-config
  ..on 'output', logger.log
  ..on 'sync-failed', -> logger.log name: 'exo-sync', text: "Some services failed to sync"
  ..on 'sync-success', -> logger.log name: 'exo-sync', text: 'Sync successful'
  ..start-syncing!
