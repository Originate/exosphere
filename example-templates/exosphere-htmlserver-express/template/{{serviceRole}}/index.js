// This is the main server file.
//
// It parses the command line and instantiates the two servers for this app:
const {cyan, dim, green, red} = require('chalk')
const ExoRelay = require('exorelay');
const fs = require('fs');
const yaml = require('js-yaml');
const N = require('nitroglycerin');
const {name, version} = require('./package.json')
const path = require('path');
const WebServer = require('./web-server')
const port = process.env.PORT || 3000


function startExorelay (done) {
  global.exorelay = new ExoRelay({
    role: process.env.ROLE,
    exocomHost: process.env.EXOCOM_HOST,
    exocomPort: process.env.EXOCOM_PORT
  })
  global.exorelay.on('error', (err) => { console.log(red(err)) })
  global.exorelay.on('online', (port) => {
    console.log(`${green('ExoRelay')} for '${process.env.ROLE}' online at port ${cyan(port)}`)
    done()
  })
  global.exorelay.connect()
}


function startWebServer (done) {
  const webServer = new WebServer;
  webServer.on('error', (err) => { console.log(red(err)) })
  webServer.on('listening', () => {
    console.log(`${green('HTML server')} online at port ${cyan(webServer.port())}`)
    done()
  })
  webServer.listen(port)
}


startExorelay( N( () => {
  startWebServer( N( () => {
    serviceConfig = yaml.safeLoad(fs.readFileSync(path.join(process.cwd(), 'service.yml'), 'utf8'))
    console.log(green(serviceConfig.startup['online-text']))
  }))
}))
