require! {
  'uuid' : uuid
  'wait' : {wait-until}
  'ws' : {Server: WebSocketServer}
}
debug = require('debug')('exocom-mock')


# ASends and receives ZMQ messages in tests
class MockExoCom

  ->
    @server = null
    @port = null
    @service-sockets = {}
    @received-messages = []
    @receive-callback = null


  close: (done) ~>
    | @server  =>  @server.close done
    | _        =>  done!


  # returns whether this server knows about an instance of the service with the given name
  knows-service: (name) ->
    !!@service-sockets[name]


  listen: (+@port, done) ~>
    @server = new WebSocketServer {@port}, done
      ..on 'error', @_on-server-error
      ..on 'connection', @_on-server-connection


  on-receive: (@receive-callback) ~>
    if @received-messages.length
      @receive-callback!


  register-service: ({name, websocket}) ~>
    @service-sockets[name] = websocket


  reset: ~>
    @received-messages = []


  send: ({service, name, payload, message-id, response-to}) ~>
    | !@service-sockets[service]  =>  throw new Error "unknown service: '#{service}'"

    @received-messages = []
    request-data =
      name: name
      payload: payload
      id: message-id or uuid.v1!
    request-data.response-to = response-to if response-to
    @service-sockets[service].send JSON.stringify request-data


  wait-until-knows-service: (name, done) ->
    wait-until (~> @knows-service name), 1, done


  _on-message: (data) ~>
    @received-messages.push data
    @receive-callback?!


  _on-server-connection: (websocket) ~>
    websocket.on 'message', @_on-socket-message(_, websocket)


  _on-server-error: (error) ->
    console.log error


  _on-socket-message: (message, websocket) ~>
    switch (request-data = JSON.parse message).name
    | 'exocom.register-service'  =>  @register-service {name: request-data.sender, websocket}
    @_on-message request-data



module.exports = MockExoCom
