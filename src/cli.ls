require! {
  './app-runner' : AppRunner
  'chalk' : {cyan, dim, green}
  'js-yaml' : yaml
  './logger' : Logger
  'fs'
  '../package.json' : {version}
}

console.log dim "Exosphere SDK #{version}\n"

app-config = yaml.safeLoad fs.readFileSync('application.yml', 'utf8')

console.log "Running #{green app-config.name} #{cyan app-config.version}\n"

logger = new Logger

app-runner = new AppRunner
  ..on 'exocomm-online', (port) ->
    logger.log name: 'exocomm', text: "online at port #{port}"

app-runner
  ..start-exocomm app-config.development['exocomm-port']
  ..on 'exocomm-online', ->
    for service of app-config.development.services
      app-runner.start-service service, app-config.development.services[service]
    console.log 'all systems go'
  ..on 'output', (data) ->
    logger.log data
