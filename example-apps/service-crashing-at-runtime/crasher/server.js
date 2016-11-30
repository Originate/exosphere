http = require('http')


requestHandler = function(req, res) {
  process.exit(1)
}


http.createServer(requestHandler)
    .listen(4000,
            '127.0.0.1',
            function() { console.log('web server running at port 4000') })
