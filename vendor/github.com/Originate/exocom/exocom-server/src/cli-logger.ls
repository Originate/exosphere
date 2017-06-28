require! {
  'chalk' : {cyan, dim, green, red, magenta, yellow}
}


# Logs events in a human-friedly format to the CLI
class CliLogger

  # logs the given critical error message
  error: (err) ~>
    @log red "Error: #{err}"
    process.exit 1


  # writes information about the application that comes before the actual application output
  header: (text) ->
    @log dim text


  # writes the given text to stdout with newline
  log: (text) ~>
    console.log text


  # logs that messages were sent over the bus
  messages: ({messages, receivers}) ~>
    for message in messages
      response-time = ''
      if message.response-to
        response-time = "  (#{(message.response-time * 1e-6).to-fixed 2} ms)"
      if message.name is message.original-name
        @log "#{message.sender}  --[ #{message.name} ]->  #{receivers.join ' and '}#{response-time}"
      else
        @log "#{message.sender}  --[ #{message.original-name} ]-[ #{message.name} ]->  #{receivers.join ' and '}#{response-time}"
      @log message.payload


  routing-setup: ~>
    @log 'receiving routing setup:'
    for command, routing of exocom.client-registry.routes
      @write "  --[ #{command} ]-> "
      text = for receiver in routing.receivers
        "#{receiver.client-name}"
      @write "#{text.join ' + '}\n"


  # logs the given warning message
  warning: (warning) ~>
    @log yellow "Warning: #{warning}"


  # writes the given text to stdout without newline
  write: (text) ~>
    process.stdout.write text



module.exports = CliLogger
