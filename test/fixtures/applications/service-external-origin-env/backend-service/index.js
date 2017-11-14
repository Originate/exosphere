require('http')
  .createServer((req, res) => res.end('Hello world!'))
  .listen(3000, () => console.log('Listening on port 3000'))
