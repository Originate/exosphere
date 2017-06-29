require! {
  './client-registry/client-registry' : ClientRegistry
  'events' : {EventEmitter}
  './http-subsystem/http-subsystem' : HttpSubsystem
  './message-cache/message-cache' : MessageCache
  'nanoseconds'
  'process'
  'rails-delegate' : {delegate, delegate-event}
  './websocket-subsystem/websocket-subsystem' : WebSocketSubsystem
}
debug = require('debug')('exocom')


class ExoCom extends EventEmitter

  ({service-routes, @logger} = {}) ->

    @client-registry = new ClientRegistry {service-routes}

    @http-subsystem = new HttpSubsystem {exocom: @, @logger}
      ..on 'online', (port) ~> @emit 'http-online', port
      ..on 'config-request', (response-stream) ~> @handle-config-request response-stream

    @message-cache = new MessageCache!

    @websocket = new WebSocketSubsystem {exocom: @, @logger}
      ..on 'online', (port) ~> @emit 'websockets-online', port
      ..on 'deregister-client', (client-name) ~> @deregister-client client-name
      ..on 'message', (message) ~> @handle-incoming-message message


  close: ->
    @http-subsystem.close!
    @websocket.close!


  # deregisters a service instance that went offline
  deregister-client: (client-name) ~>
    @client-registry.deregister-client client-name
    @websocket.deregister-client client-name


  # returns the current configuration of this ExoCom instance
  get-config: ~>
    {
      clients: @client-registry.clients
      routes: @client-registry.routing
    }


  # returns whether the given message contains an invalid sender
  has-invalid-sender: (message) ->
    !@client-registry.can-send message.sender, message.name


  # bind to the given port to send socket messages
  listen: (port) ->
    express-server = @http-subsystem.listen port
    @websocket.listen port, express-server
    debug "Listening at port #{port}"


  # called when the HTTP subsystem emits a request for configuration data
  handle-config-request: (response-stream) ->
    @http-subsystem.send-configuration {configuration: @get-config!, response-stream}


  # called when a new message from a service instance arrives
  handle-incoming-message: ({message, websocket}) ->
    | message.name is \exocom.register-service  =>  @register-client {message, websocket}
    | @has-invalid-sender message               =>  @logger.error "Service '#{message.sender}' is not allowed to broadcast the message '#{message.name}'"
    | otherwise                                 =>  @send-message message


  # called when a new service instance registers with this Exocom instance
  register-client: ({message, websocket}) ->
    @client-registry.register-client message.payload
    @websocket.register-client {client-name: message.payload.client-name, websocket}


  # sends the given message to all subscribers of it.
  send-message: (message-data) ~>
    # convert the outgoing message name from its internal version to the public version
    sender = @client-registry.clients[message-data.sender]
    public-message-name = @client-registry.outgoing-message-name message-data.name, sender
    message-data.original-name = message-data.name
    message-data.name = public-message-name
    message-data.timestamp = nanoseconds process.hrtime!
    # determine the subscribers
    subscribers = @client-registry.subscribers-for public-message-name
    return @logger.warning "No receivers for message '#{message-data.name}' registered" unless subscribers
    subscriber-names = [subscriber.client-name for subscriber in subscribers]

    # calculate a message's response time if it is a reply
    if message-data.response-to
      original-timestamp  = @message-cache.get-original-timestamp message-data.id
      message-data.response-time = message-data.timestamp - original-timestamp
    else
      @message-cache.push message-data.id, message-data.timestamp

    # send the message to the subscribers
    debug "sending '#{message-data.name}' to #{subscriber-names}"
    sent-messages = @websocket.send-message-to-services message-data, subscribers
    @logger.messages messages: sent-messages, receivers: subscriber-names

    'success'



module.exports = ExoCom
