require! {
  './app-runner' : AppRunner
  'chalk' : {bold, cyan, dim, green, red}
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
  ..on 'exocom-online', (port) -> logger.log name: 'exocom', text: "online at port #{port}"
  ..on 'service-online', (name) -> logger.log name: 'exorun', text: "'#{name}' is running using exorelay port #{app-runner.port-for name}"
  ..on 'routing-setup', ->
    logger.log name: 'exocom', text: 'received routing setup'
    for command, routing of app-runner.exocom.client-registry.routes
      text = "  [ #{bold command} ]  -->  "
      receivers = for receiver in routing.receivers
        "#{bold receiver.name} (#{receiver.host}:#{receiver.port})"
      text += receivers.join ' & '
      logger.log name: 'exocom', text: text

  ..on 'message', ({sender, message, receivers}) -> logger.log name: 'exocom', text: "#{sender}  --[ #{message} ]->  #{receivers.join ' and '}"
  ..start-exocom!
  ..start-services!
  ..on 'all-services-online', ->
    logger.log name: 'exorun', text: 'all services online'
    app-runner
      ..on 'routing-done', -> logger.log name: 'exorun', text: "application ready"
      ..send-service-configuration!
