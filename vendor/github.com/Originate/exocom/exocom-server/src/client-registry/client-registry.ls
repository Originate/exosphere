require! {
  'remove-value'
  'require-yaml'
  './subscription-manager' : SubscriptionManager
}


class ClientRegistry

  ({service-routes = '{}'} = {}) ->

    # List of messages that are received by the applications services
    #
    # the format is:
    # {
    #   'role 1':
    #     receives: ['message 1 name', 'message 2 name']
    #     sends: ['message 3 name', 'message 4 name']
    #     internal-namespace: 'my internal namespace'
    #   'role 2':
    #     ...
    @routing = @_parse-service-routes service-routes

    # The main list of clients that are currently registered
    #
    # The format is:
    # {
    #   'client 1 name':
    #     client-name: ...
    #     service-type: ...
    #     namespace: ...
    #   'client 2 name':
    #     ...
    @clients = {}

    @subscriptions = new SubscriptionManager @routing


  # returns whether the given sender is allowed to send messages with the given name
  can-send: (sender-name, message-name) ->
    @routing[sender-name].sends |> (.includes message-name)


  # deregisters a service instance that went offline
  deregister-client: (client-name) ->
    @subscriptions.remove client-name
    delete @clients[client-name]


  # Returns the external name for the given message sent by the given service,
  # i.e. how the sent message should appear to the other services.
  outgoing-message-name: (message-name, service) ->
    message-parts = message-name.split '.'
    switch
    | message-parts.length is 1                       =>  message-name
    | message-parts[0] is service.internal-namespace  =>  "#{service.service-type}.#{message-parts[1]}"
    | otherwise                                       =>  message-name


  # registers the given service instance that just came online
  register-client: (client) ->
    @clients[client.client-name] =
      client-name: client.client-name
      service-type: client.client-name
      internal-namespace: @routing[client.client-name].internal-namespace

    @subscriptions.add-all {client.client-name, service-type: client.client-name}


  # Returns the clients that are subscribed to the given message
  subscribers-for: (message-name) ->
    @subscriptions.subscribers-for message-name



  _parse-service-routes: (service-routes) ->
    result = {}
    for service-route in JSON.parse(service-routes)
      result[service-route.role] =
        receives: service-route.receives
        sends: service-route.sends or []
        internal-namespace: service-route.namespace
    result



module.exports = ClientRegistry
