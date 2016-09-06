require! {
  './app-setup' : AppSetup
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  'exosphere-shared' : {Logger}
  'fs'
}

app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
console.log "Setting up #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger Object.keys(app-config.services)
app-setup = new AppSetup app-config
  ..on 'start', (name) -> logger.log name: 'exo-setup', text: "starting setup of '#{name}'"
  ..on 'error', (name, exit-code) ->
    console.log red "\nsetup of '#{name}' failed with exit code #{exit-code}"
    process.exit exit-code
  ..on 'output', (data) -> data.text = data.text.replace('\n', '') ; logger.log data
  ..on 'finished', (name) -> logger.log name: 'exo-setup', text: "setup of '#{name}' finished"
  ..on 'setup-complete', -> logger.log name: 'exo-setup', text: 'setup complete'
  ..start-setup!
