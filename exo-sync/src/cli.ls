require! {
  'chalk' : {cyan}
  'fs'
  'js-yaml': yaml
  './app-syncer' : AppSyncer
  '../../exosphere-shared' : {Logger}
}

module.exports = ->

  if process.argv[2] is \help
    return help!

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf-8')
  logger = new Logger Object.keys(app-config.services ? {})
    ..log role: 'exo-sync', text: "Syncing application #{cyan app-config.name}"
  syncer = new AppSyncer {app-config, logger}
    ..start-syncing!

function help
  help-message =
    """
    Usage: #{cyan "exo sync"}

    Syncs all services of an Exosphere application with their git repositories.
    This command must be run in the root directory of the application.
    """
  console.log help-message
