require! {
  '../message-translator/message-translator' : MessageTranslator
}

# manages which client is subscribed to which message on the bus
class SubscriptionManager


  (@routing) ->
    @message-translator = new MessageTranslator

    # List of clients that are subscribed to the given message
    #
    # The format is:
    # {
    #   'message 1 name': [
    #     * name: 'client 1 name'
    #       internal-namespace: 'my namespace'
    #     * name: ...
    #   ],
    #   'message 2 name':
    #     ...
    @subscribers = {}


  # Adds the given client to the subscription list for the given message
  add: ({internal-message-name, client-name}) ->
    public-message-name = @message-translator.public-message-name {internal-message-name, client-name, internal-namespace: @routing[client-name].internal-namespace}
    (@subscribers[public-message-name] or= []).push do
      client-name: client-name
      internal-namespace: @routing[client-name].internal-namespace


  # adds subscriptions for the client with the given name
  add-all: ({client-name, service-type}) ->
    for internal-message-name in @routing[service-type].receives or {}
      @add {internal-message-name, client-name}


  remove: (client-name) ->
    for internal-message-name in (@routing[client-name].receives or {})
      public-message-name = @message-translator.public-message-name {internal-message-name, client-name, internal-namespace: @routing[client-name].internal-namespace}
      # TODO: this is broken, make this remove only the client
      delete @subscribers[public-message-name]


  subscribers-for: (message-name) ->
    @subscribers[message-name]



module.exports = SubscriptionManager
