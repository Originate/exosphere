require! {
  'http'
  'exorelay' : ExoRelay
}

exorelay = new ExoRelay exocom-host: process.env.EXOCOM_HOST, exocom-port: 8000, service-name: 'web'
  ..on 'online', (port) -> console.log "web service exorelay online at port #{port}"
  ..on 'error', (err) -> console.log "web service exorelay encountered error: #{err}"
  ..listen 8002


request-handler = (req, res) ->
  exorelay.send 'users.list', {}, (_, {outcome}) ->
    console.log "'web' service received message '#{outcome}'"
    res.writeHead 200, 'Content-Type': 'text/plain'
    res.end 'test web server\n'


http.create-server request-handler
  ..listen 4000 ->
    console.log "web server running at port 4000"
