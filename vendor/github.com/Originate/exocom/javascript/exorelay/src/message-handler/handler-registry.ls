require! {
  'chalk' : {bold, red}
  'debug'
  'events' : {EventEmitter}
}


# A registry for message handlers of a particular type
class HandlerRegistry extends EventEmitter

  (debug-name) ->

    # handler functions for incoming requests
    @handlers = {}

    @debug = debug "exorelay:#{debug-name}"


  # Returns the handler for the request with the given id,
  # or undefined if not found.
  get-handler: (message-id) ->
    @handlers[message-id]


  # Tries to handle the given incoming command.
  # Returns whether it was able to do this.
  handle-command: ({message-name, payload}, methods) ->
    if handler = @get-handler message-name
      debug "handling message '#{message-name}'"
      try
        handler payload, methods
      catch
        console.log "\n#{red bold e.message}\n"
        console.log e.stack
        throw e
    !!handler


  # Tries to handle the incoming reply.
  # Returns whether it was able to do that.
  handle-reply: ({message-name, response-to, payload}) ->
    if handler = @get-handler response-to
      debug "handling message '#{message-name}' in response to '#{response-to}'"
      try
        handler payload, outcome: message-name
      catch
        console.log "\n#{red bold e.message}\n"
        console.log e.stack
        throw e
    !!handler


  # Returns whether this RequestManager has a handler for
  # the request with the given id registered.
  has-handler: (message-id) ->
    typeof @get-handler(message-id) is 'function'


  register-handler: (message-id, handler) ->
    | !message-id                      =>  return @emit 'error', new Error 'No message id provided'
    | typeof message-id isnt 'string'  =>  return @emit 'error', new Error 'Message ids must be strings'
    | !handler                         =>  return @emit 'error', new Error 'No message handler provided'
    | typeof handler isnt 'function'   =>  return @emit 'error', new Error 'Message handler must be a function'
    | @has-handler message-id          =>  return @emit 'error', new Error "There is already a handler for message '#{message-id}'"

    @debug "registering handler for id '#{message-id}'"
    @handlers[message-id] = handler


  register-handlers: (handlers) ->
    for message-id, handler of handlers
      @register-handler message-id, handler



module.exports = HandlerRegistry
