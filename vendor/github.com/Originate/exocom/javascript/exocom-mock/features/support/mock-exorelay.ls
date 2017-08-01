require! {
  'ws' : WebSocket
}
debug = require('debug')('mock-websocket-endpoint')


# A programmable Exorelay mock for testing MockExocom
class MockExorelay

  (@name) ->
    @received-messages = []
    @socket = null


  connect: ({@exocom-port, registration-message, registration-payload}, done) ~>
    @socket = new WebSocket "ws://localhost:#{@exocom-port}"
      ..on 'error', @_on-socket-error
      ..on 'open', ~> @register {message-name: registration-message, payload: registration-payload} ; done!
      ..on 'message', @_on-socket-message


  close: ~>
    @socket?.close!


  register: ({message-name = 'exocom.register-service', payload}) ->
    @send do
      name: message-name
      sender: @name
      payload: payload


  send: (request-data) ~>
    @socket.send JSON.stringify request-data



  _on-socket-error: (error) ->
    console.log error


  _on-socket-message: (data) ~>
    @received-messages.push JSON.parse data.to-string!



module.exports = MockExorelay
