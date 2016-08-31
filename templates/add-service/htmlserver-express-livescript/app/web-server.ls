require! {
  'events' : {EventEmitter}
  'express'
  'exprestive'
  'jade'
  'http'
  'method-override'
  'path'
  'rails-delegate' : {delegate, delegate-event}
  'morgan' : logger
  'cookie-parser'
  'body-parser'
}


class WebServer extends EventEmitter

  ->
    @app = express!
      ..use methodOverride '_method'

    # view engine setup
    @app.set 'views', path.join __dirname, 'views'
      ..set 'view engine', \jade

      ..use logger \dev
      ..use express.static path.join __dirname, 'public'
      ..use require '../webpack/middleware'
      ..use bodyParser.json!
      ..use bodyParser.urlencoded extended: false
      ..use cookieParser!

      ..use exprestive dependencies: global.exorelay

      ..use (req, res, next) ->   # route not found
        err = new Error 'Not Found'
        err.status = 404
        next err

    @server = http.create-server @app

    delegate \listen, from: @, to: @server
    delegate-event \listening \error, from: @server, to: @


  port: ->
    @server.address!port



module.exports = WebServer
