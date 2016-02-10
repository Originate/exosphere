require! {
  'docopt' : {docopt}
  'http'
}


doc = """
Test server

Usage:
  start --html-port=<html-port> --exorelay-port=<exorelay-port> --exocomm-port=<exocomm-port>
"""
options = docopt doc, help: no


request-handler = (req, res) ->
  res.writeHead 200, 'Content-Type': 'text/plain'
  res.end 'test web server\n'

http.create-server request-handler
    .listen options['--html-port'], '127.0.0.1', ->
      console.log "Server running at port #{options['--html-port']}"
