require! {
  './app-installer' : AppInstaller
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  '../../logger' : Logger
  'fs'
  '../../../package.json' : {version}
}

console.log dim "Exosphere SDK #{version}\n"
app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
console.log "Installing #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger Object.keys(app-config.services)
app-installer = new AppInstaller app-config
  ..on 'start', (name) -> logger.log name: 'exo-install', text: "starting setup of '#{name}'"
  ..on 'error', (err) -> console.log red error
  ..on 'output', (data) -> data.text = data.text.replace('\n', '') ; logger.log data
  ..on 'finished', (name) -> logger.log name: 'exo-install', text: "setup of '#{name}' finished"
  ..on 'installation-complete', -> logger.log name: 'exo-install', text: 'installation complete'
  ..start-installation!
