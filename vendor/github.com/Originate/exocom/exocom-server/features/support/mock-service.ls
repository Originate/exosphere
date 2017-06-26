require! {
  'ws' : WebSocket
}
debug = require('debug')('mock-service')


# A class to encapsulate the functionality of a service in the ExoSphere
# that can send, receive, and track messages using WebSockets.
class MockService

  ({@port, @client-name, @namespace} = {}) ->
    @received-messages = []


  close: ~>
    | @closed => return
    @socket?.close!
    @closed = yes


  connect: ({payload}, done) ~>
    payload ?= {@client-name}
    @socket = new WebSocket "ws://localhost:#{@port}/services"
      ..on 'message', @_on-socket-message
      ..on 'error', @_on-socket-error
      ..on 'open', ~>
        @send do
          name: 'exocom.register-service'
          sender: @client-name
          payload: payload
          id: '123'
        done!


  send: (request-data) ~>
    @socket.send JSON.stringify request-data


  _on-socket-error: (error) ~>
    console.log error


  _on-socket-message: (data) ~>
    @received-messages.unshift(JSON.parse data.to-string!)



module.exports = MockService
