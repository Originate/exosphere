const EventEmitter = require('events').EventEmitter
const express = require('express')
const exprestive = require('exprestive')
const jade = require('jade')
const http = require('http')
const methodOverride = require('method-override')
const path = require('path')
const ref$ = require('rails-delegate'), delegate = ref$.delegate, delegateEvent = ref$.delegateEvent
const logger = require('morgan')
const cookieParser = require('cookie-parser')
const bodyParser = require('body-parser')

class WebServer extends EventEmitter {

  constructor() {
    super()
    this.app = express()
    this.app.use(methodOverride('_method'))

    // view engine setup
    this.app.set('views', path.join(__dirname, 'views'))
    this.app.set('view engine', 'jade')
    this.app.use(logger('dev'))
    this.app.use(express['static'](path.join(__dirname, 'public')))
    this.app.use(require('../webpack/middleware'))
    this.app.use(bodyParser.json())
    this.app.use(bodyParser.urlencoded({ extended: false }))
    this.app.use(cookieParser())

    this.app.use(exprestive({ dependencies: global.exorelay }))

    this.app.use((req, res, next) => {   // route not found
      const err = new Error('Not Found')
      err.status = 404
      next(err)
    })

    this.server = http.createServer(this.app)

    delegate('listen', { from: this, to: this.server })
    delegateEvent('listening', 'error', { from: this.server, to: this })
  }


  port() {
    return this.server.address().port
  }

}



module.exports = WebServer
