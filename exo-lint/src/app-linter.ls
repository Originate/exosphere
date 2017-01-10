require! {
  'async'
  'events' : {EventEmitter}
  'fs'
  'path'
  'prelude-ls' : {difference, find, filter, reject, each}
  'js-yaml' : yaml
}

class AppLinter extends EventEmitter

  ({@app-config, @logger}) ->


  start: ->
    {sent-messages, received-messages} = @aggregate-messages!
    @lint-messages sent-messages, received-messages


  lint-messages: (sent, received) ->
    not-received = difference Object.keys(sent), Object.keys(received)
    not-sent = difference Object.keys(received), Object.keys(sent)

    if not-received.length is 0 and not-sent.length is 0
      return @logger.log name: 'exo-lint', text: 'Lint passed'

    if not-received.length
      @logger.log name: 'exo lint', text: "The following messages are sent but not received:"
      for msg in not-received
        @logger.log name: sent[msg], text: msg
    if not-sent.length
      @logger.log name: 'exo lint', text: "The following messages are received but not sent:"
      for msg in not-sent
        @logger.log name: received[msg], text: msg


  aggregate-messages: ->
    sent-messages = {}
    received-messages = {}
    for protection-level of @app-config.services
      for service-role, service-data of @app-config.services[protection-level]
        service-config = @get-config service-data
        for message in service-config.messages.sends or []
          (sent-messages[message] or= []).push service-role
        for message in service-config.messages.receives or []
          (received-messages[message] or= []).push service-role

    {sent-messages, received-messages}


  get-config: (service-data) ->
    yaml.safe-load fs.read-file-sync(path.join(process.cwd!, service-data.location, 'service.yml'), 'utf8')


module.exports = AppLinter
