http = require('http')


requestHandler = function(req, res) {
  res.writeHead(200, {'Content-Type': 'text/plain'})
  res.end('simple example app\n')
}


http.createServer(requestHandler)
    .listen(4000,
            '127.0.0.1',
            function() { console.log('web service running at port 4000') })
