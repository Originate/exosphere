require! {
  'async'
  'events' : {EventEmitter}
  'fs'
  'path'
  'prelude-ls' : {difference, find, filter, reject, each}
  'js-yaml' : yaml
}

class AppLinter extends EventEmitter

  (@app-config) ->


  start: ->
    {sent-messages, received-messages} = @aggregate-messages!
    @emit 'reset colors', Object.keys @app-config.services
    @lint-messages sent-messages, received-messages


  lint-messages: (sent, received) ->
    not-received = difference Object.keys(sent), Object.keys(received)
    not-sent = difference Object.keys(received), Object.keys(sent)

    if not-received.length is 0 and not-sent.length is 0
      return @emit 'lint success'

    if not-received.length
      @emit 'output', {name: 'exo lint', text: "The following messages are sent but not received:"}
      for msg in not-received
        @emit 'output', {name: sent[msg], text: msg}
    if not-sent.length
      @emit 'output', {name: 'exo lint', text: "The following messages are received but not sent:"}
      for msg in not-sent
        @emit 'output', {name: received[msg], text: msg}


  aggregate-messages: ->
    sent-messages = {}
    received-messages = {}
    for service-name of @app-config.services
      service-config = @get-config service-name
      for message in service-config.messages.sends or []
        (sent-messages[message] or= []).push service-config.name
      for message in service-config.messages.receives or []
        (received-messages[message] or= []).push service-config.name

    {sent-messages, received-messages}


  get-config: (service-name) ->
    service-root = path.join process.cwd!, @app-config.services[service-name].location
    yaml.safe-load fs.read-file-sync(path.join(service-root, 'service.yml'), 'utf8')



module.exports = AppLinter
