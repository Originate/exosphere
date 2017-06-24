require! {
  'chalk' : {cyan}
  'events' : {EventEmitter}
  '../message-cache/message-cache' : MessageCache
  '../message-translator/message-translator' : MessageTranslator
  'ws' : {Server: WebSocketServer}
}
debug = require('debug')('exocom:websocket-subsystem')


# The web sockets endpoint that listens/sends messages from/to services
#
# Emits these events:
# - online: when it is online and ready to go
# - error: for critical issues
# - warn: for non-critical issues
class WebSocketSubsystem extends EventEmitter

  ({@logger} = {}) ->
    @message-translator = new MessageTranslator
    @server = null
    @port = null

    # Stores the service names and respective web socket
    #
    # format:
    # {
    #   'client 1 name': websocket,
    #   'client 2 name': websocket,
    #   ...
    # }
    @sockets = {}


  close: ->
    | !@server  =>  return
    debug 'websockets going offline'
    @server.close!


  deregister-client: (client-name) ->
    @sockets[client-name]?.close!
    delete @sockets[client-name]


  # Listens at the given port
  # by hooking into the given http server instance
  listen: (@port, server) ->
    @server = new WebSocketServer {server, path: '/services'}
      ..on 'connection', @on-connection
      ..on 'listening', ~>
        @logger.log "ExoCom WebSocket listener online at port #{cyan port}"
        @emit 'online', @port
      ..on 'error', (err) ~> @logger.error err


  # called when a new service instance connects
  on-connection: (websocket) ~>
    websocket.on 'message', (message-text) ~>
      message = JSON.parse message-text
      @_log-received message
      @emit 'message', {message, websocket}


  # Registers the given websocket as a connection
  # to an instance of the service with the given name
  register-client: ({client-name, websocket}) ->
    @sockets[client-name] = websocket
      ..on 'close', ~> @emit 'deregister-client', client-name


  send-message-to-service: (message, service) ->
    internal-message-name = @message-translator.internal-message-name message.name, for: service
    request-data =
      name: internal-message-name
      id: message.id
      payload: message.payload
      timestamp: message.timestamp
    if message.response-to
      request-data.response-time = message.response-time
      request-data.response-to = message.response-to
    @_log-sending message, service
    @sockets[service.client-name].send JSON.stringify request-data
    result = {[key, value] for key, value of message}
    result.name = internal-message-name
    result


  send-message-to-services: (message-data, services) ->
    for service in services
      @send-message-to-service message-data, service


  _log-received: (message) ->
    | message.response-to  =>  debug "received '#{message.name}' with id '#{message.id}' in response to '#{message.response-to}'"
    | _                    =>  debug "received '#{message.name}' with id '#{message.id}'"


  _log-sending: (message, service) ->
    | message.response-to  =>  debug "sending '#{message.name}' with id '#{message.id}' in response to '#{message.response-to}' to '#{service.name}'"
    | _                    =>  debug "sending '#{message.name}' with id '#{message.id}' to '#{service.name}'"



module.exports = WebSocketSubsystem
