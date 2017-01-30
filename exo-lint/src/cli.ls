require! {
  './app-linter' : AppLinter
  'chalk' : {cyan, blue}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'fs'
}

module.exports = ->

  if process.argv[2] is \help
    return help!

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  console.log "Running linter for #{cyan app-config.name}\n"
  logger = new Logger Object.keys(app-config.services)
  app-linter = new AppLinter {app-config, logger}
    ..start!

function help
  help-text = """
  Usage: #{cyan "exo lint"}

  Runs a linter for the current application.
  This linter checks that all messages indicated to be sent by services of the application are received by one or more other services in the application and vice versa.
  This command must be called in the root directory of the application.
  """
  console.log help-text
