require! {
  './app-runner' : AppRunner
  'chalk' : {cyan, dim, green, red}
  'js-yaml' : yaml
  'fs'
  '../package.json' : {version}
}

console.log dim "Exosphere SDK #{version}\n"

try
  app-config = yaml.safeLoad fs.readFileSync('application.yml', 'utf8')
catch
  console.log(e)

console.log "Running #{green app-config.name} #{cyan app-config.version}\n"

console.log app-config

app-runner = new AppRunner
  ..on 'exocomm-online', (port) -> console.log "Exocomm online at port #{port}"

app-runner
  ..start-exocomm app-config.development['exocomm-port']
  ..on 'exocomm-online', ->
    console.log 'all systems go'
