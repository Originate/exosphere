require! {
  'http'
  'exorelay' : ExoRelay
}

exorelay = new ExoRelay exocom-host: process.env.EXOCOM_HOST, exocom-port: process.env.EXOCOM_PORT, role: 'web'
  ..on 'online', -> console.log "web service exorelay online"
  ..on 'error', (err) -> console.log "web service exorelay encountered error: #{err}"
  ..connect!


request-handler = (req, res) ->
  exorelay.send 'users.list', {}, (_, {outcome}) ->
    console.log "'web' service received message '#{outcome}'"
    res.writeHead 200, 'Content-Type': 'text/plain'
    res.end 'test web server\n'


http.create-server request-handler
  ..listen 4000 ->
    console.log "web server running at port 4000"
