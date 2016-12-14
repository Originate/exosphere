# This is the main server file.
#
# It parses the command line and instantiates the two servers for this app:
require! {
  'chalk' : {cyan, dim, green, red}
  'exorelay' : ExoRelay
  'fs'
  'js-yaml': yaml
  'nitroglycerin' : N
  '../package.json' : {name, version}
  'path'
  './web-server' : WebServer
}


start-exorelay = (done) ->
  global.exorelay = new ExoRelay service-name: process.env.SERVICE_NAME, exocom-host: process.env.EXOCOM_HOST, exocom-port: process.env.EXOCOM_PORT
    ..on 'error', (err) -> console.log red err
    ..on 'online', (port) ->
      console.log "#{green 'ExoRelay'} for '#{process.env.SERVICE_NAME}' online at port #{cyan port}"
      done!
    ..connect!


start-web-server = (done) ->
  web-server = new WebServer
    ..on 'error', (err) -> console.log red err
    ..on 'listening', ->
      console.log "#{green 'HTML server'} online at port #{cyan web-server.port!}"
      done!
    ..listen 3000


start-exorelay N ->
  start-web-server N ->
    service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, 'service.yml'), 'utf8')
    console.log green service-config.startup[\online-text]
