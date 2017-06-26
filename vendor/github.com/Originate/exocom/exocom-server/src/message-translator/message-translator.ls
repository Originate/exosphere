class MessageTranslator

  # Returns the message name to which the given service would have to subscribe
  # if it wanted to receive the given message expressed in its internal form.
  #
  # Example:
  # - service "tweets" has internal namespace "text-snippets"
  # - it only knows the "text-snippets.create" message
  # - the external message name that it has to subscribe to is "tweets.create"
  public-message-name: ({internal-message-name, client-name, internal-namespace}) ->
    | !internal-namespace                                  =>  internal-message-name
    | (internal-message-name.split '.').length is 1        =>  internal-message-name
    | (internal-message-name.split '.')[0] is client-name  =>  internal-message-name
    | otherwise                                            =>  "#{client-name}.#{(internal-message-name.split '.')[1]}"


  # Translates outgoing message into one that the receiving service will understand
  internal-message-name: (message-name, {for: service}) ->
    | !service.internal-namespace                                =>  message-name
    | (message-name.split '.').length is 1                       =>  message-name
    | (message-name.split '.')[0] is service.internal-namespace  =>  message-name
    | otherwise                                                  => "#{service.internal-namespace}.#{(message-name.split '.')[1]}"


module.exports = MessageTranslator
