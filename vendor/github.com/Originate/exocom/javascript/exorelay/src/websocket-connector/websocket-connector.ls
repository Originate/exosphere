require! {
  'events' : {EventEmitter}
  'uuid' : uuid
  'wait' : {wait}
  'ws' : WebSocket
}
debug = require('debug')('exorelay:websocket-listener')


# The WebSocket endpoint which the Exosphere environment can
# send messages and receive messages from.
#
# Emits these events:
# - online
# - offline
# - error
class WebSocketConnector extends EventEmitter

  ({@exocom-host, @role, @exocom-port} = {}) ->
    @exocom-port = +@exocom-port

    # Contains the id of the most recently sent request (for testing)
    @last-sent-id = null


  # Closes the port that ExoRelay is communicating on
  close: ->
    return unless @socket
    debug "no longer connected at 'ws://#{@exocom-host}/#{@exocom-port}'"
    @socket?.close!
    @emit 'offline'


  # Returns a method that sends a reply to the message with the given request
  reply-method-for: (id) ->
    | !id  =>  return @emit 'error', new Error 'WebSocketConnector.replyMethodFor needs an id'

    (message-name, payload = {}) ~>
      @send message-name, payload, response-to: id


  send: (message-name, payload, options = {}) ->
    | !message-name                      =>  return @emit 'error', new Error 'ExoRelay#send cannot send empty messages'
    | typeof message-name isnt 'string'  =>  return @emit 'error', new Error 'ExoRelay#send can only send string messages'
    | typeof payload is 'function'       =>  return @emit 'error', new Error 'ExoRelay#send cannot send functions as payload'

    @_log-sending message-name, options
    request-data =
      name: message-name
      sender: @role
      id: uuid.v1!
    request-data.payload = payload if payload?
    request-data.response-to = options.response-to if options.response-to
    @socket.send JSON.stringify request-data
    @last-sent-id = request-data.id


  connect: ~>
    @socket = new WebSocket "ws://#{@exocom-host}:#{@exocom-port}/services"
      ..on 'open', @_on-socket-open
      ..on 'message', @_on-socket-message
      ..on 'error', @_on-socket-error


  _on-socket-open: ~>
    @emit 'online'


  _on-socket-error: (error) ~>
    | error.errno is 'EADDRINUSE'   =>  @emit 'error', "port #{@exocom-port} is already in use"
    | error.errno is 'ECONNREFUSED' =>  wait 1_000, @connect
    | otherwise                     =>  @emit 'error', error


  _on-socket-message: (data) ~>
    request-data = data |> JSON.parse |> @_parse-request
    @_log-received request-data
    switch (result = @listeners('message')[0] request-data)
      | 'success'             =>
      | 'missing message id'  =>  @emit 'error', Error 'missing message id'
      | 'unknown message'     =>  @emit 'error', Error "unknown message: '#{request-data.message-name}'"
      | _                     =>  @emit 'error', Error "unknown result code: '#{result}'"


  _log-received: ({message-name, id, response-to}) ->
    | response-to  =>  debug " received message '#{message-name}' with id '#{id}' in response to '#{response-to}'"
    | _            =>  debug "received message '#{message-name}' with id '#{id}'"


  _log-sending: (message-name, options) ->
    | options.response-to  =>  debug "sending message '#{message-name}' in response to '#{options.response-to}'"
    | _                    =>  debug "sending message '#{message-name}'"


  # Returns the relevant data from a request
  _parse-request: (req) ->
    {
      message-name: req.name
      payload: req.payload
      response-to: req.response-to
      id: req.id
    }



module.exports = WebSocketConnector
