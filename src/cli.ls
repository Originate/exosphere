require! {
  './app-runner' : AppRunner
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  'fs'
  '../package.json' : {version}
}

console.log dim "Exosphere SDK #{version}\n"

app-config = yaml.safeLoad fs.readFileSync('application.yml', 'utf8')

console.log "Running #{green app-config.name} #{cyan app-config.version}\n"

app-runner = new AppRunner
  ..on 'exocomm-online', (port) -> console.log "Exocomm online at port #{port}"

app-runner
  ..start-exocomm app-config.development['exocomm-port']
  ..on 'exocomm-online', ->
    for service of app-config.development.services
      app-runner.start-service service, app-config.development.services[service]
    console.log 'all systems go'
