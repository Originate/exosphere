require! {
  './app-runner' : AppRunner
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  '../../logger' : Logger
  'fs'
  '../../../package.json' : {version}
}

console.log dim "Exosphere SDK #{version}\n"
app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
console.log "Running #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger Object.keys(app-config.services)
app-runner = new AppRunner app-config
  ..on 'error', (err) -> console.log red error
  ..on 'output', (data) -> logger.log data
  ..on 'exocomm-online', (port) -> logger.log name: 'exocomm', text: "online at port #{port}"
  ..on 'service-online', (name) -> logger.log name: 'exorun', text: "'#{name}' is running"
  ..on 'routing-setup', -> logger.log name: 'exocomm', text: 'received routing setup'
  ..on 'command', (name, receivers) -> logger.log name: 'exocomm', text: "broadcasting '#{name}' to #{receivers.join ' and '}"
  ..start-exocomm!
  ..start-services!
  ..on 'all-services-online', ->
    logger.log name: 'exorun', text: 'all services online'
    app-runner
      ..on 'routing-done', -> logger.log name: 'exorun', text: "application ready"
      ..send-service-configuration!
