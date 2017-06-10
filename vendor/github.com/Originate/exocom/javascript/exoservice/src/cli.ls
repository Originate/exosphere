require! {
  \chalk : {cyan, dim, green, red}
  'fs'
  'js-yaml' : yaml
  'path'
  '../dist/exo-js' : ExoService
  '../package.json' : exo-js-data
}

console.log dim "Exosphere Node.js service runner #{exo-js-data.version}\n"

service-data = yaml.safe-load fs.read-file-sync(path.resolve('./service.yml'), 'utf8')
console.log "Running #{green process.env.ROLE}\n"

new ExoService parse-options!
  ..on 'online', (port) ->
    console.log dim "Ctrl-C to stop"
    console.log "online at port #{cyan port}"
  ..on 'error', (err) -> console.log red err
  ..on 'offline', -> console.log red 'SERVER CLOSED'
  ..connect!


function parse-options
  exocom-host: process.env.EXOCOM_HOST ? "localhost"
  exocom-port: process.env.EXOCOM_PORT
  role: process.env.ROLE
  internal-namespace: service-data.namespace
