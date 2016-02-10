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
new AppRunner
  ..start-exocomm app-config.development['exocomm-port']
  ..on 'exocomm-online', (port) -> logger.log name: 'exocomm', text: "online at port #{port}"
  ..on 'output', (data) -> logger.log data
  ..start-services app-config.development.services
  ..on 'service-online', (name) -> logger.log name: 'exorun', text: "'#{name}' is running"
  ..on 'all-services-online', -> logger.log name: 'exorun', text: 'all systems online'
