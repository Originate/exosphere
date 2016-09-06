# This is the main server file.
#
# It parses the command line and instantiates the two servers for this app:
require! {
  'chalk' : {cyan, dim, green, red}
  'exorelay' : ExoRelay
  'nitroglycerin' : N
  '../package.json' : {name, version}
  './web-server' : WebServer
}


start-exorelay = (done) ->
  global.exorelay = new ExoRelay service-name: process.env.SERVICE_NAME, exocom-port: process.env.EXOCOM_PORT
    ..on 'error', (err) -> console.log red err
    ..on 'online', (port) ->
      console.log "#{green 'ExoRelay'} for '#{process.env.SERVICE_NAME}' online at port #{cyan port}"
      done!
    ..listen process.env.EXORELAY_PORT


start-web-server = (done) ->
  web-server = new WebServer
    ..on 'error', (err) -> console.log red err
    ..on 'listening', ->
      console.log "#{green 'HTML server'} online at port #{cyan web-server.port!}"
      done!
    ..listen 3000


start-exorelay N ->
  start-web-server N ->
    console.log green 'HTML server is running'
