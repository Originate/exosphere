require! {
  'docopt' : {docopt}
  'http'
}


doc = """
Test server

Usage:
  exo-js --name=<name> --exorelay-port=<exorelay-port> --exocomm-port=<exocomm-port>
"""
options = docopt doc, help: no


request-handler = (req, res) ->
  res.writeHead 200, 'Content-Type': 'text/plain'
  res.end 'test web server\n'

http.create-server request-handler
  .listen 4000, '127.0.0.1', ->
    console.log "Server running at port 4000"
