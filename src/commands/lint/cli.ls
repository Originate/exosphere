require! {
  './app-linter' : AppLinter
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  '../../logger' : Logger
  'fs'
}

app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
console.log "Running linter for #{cyan app-config.name}\n"
logger = new Logger Object.keys(app-config.services)
app-linter = new AppLinter app-config
  ..on 'lint success', -> logger.log name: 'exo-lint', text: 'Lint passed'
  ..on 'output', (data) -> logger.log data
  ..on 'reset colors', (service-names) -> logger.set-colors service-names
  ..start!
