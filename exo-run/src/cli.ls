require! {
  './app-runner' : AppRunner
  'chalk' : {bold, cyan, dim, green, red}
  'fs'
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'prelude-ls' : {flatten}
  'util'
}

module.exports = ->

  if process.argv[2] is "help"
    help!
    return

  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  console.log "Running #{green app-config.name} #{cyan app-config.version}\n"
  logger = new Logger flatten [Object.keys(app-config.services[type]) for type of app-config.services]
  app-runner = new AppRunner {app-config, logger}
    ..on 'routing-setup', ->
      logger.log role: 'exocom', text: 'received routing setup'
      for command, routing of app-runner.exocom.client-registry.routes
        text = "  [ #{bold command} ]  -->  "
        receivers = for receiver in routing.receivers
          "#{bold receiver.name} (#{receiver.host}:#{receiver.port})"
        text += receivers.join ' & '
        logger.log role: 'exocom', text: text

    ..on 'message', ({messages, receivers}) ->
      message = messages[0]
      if message.name isnt message.original-name
        logger.log role: 'exocom', text: "#{bold message.sender}  --[ #{bold message.original-name} ]-[ #{bold message.name} ]->  #{bold receivers.join ' and '}"
      else
        logger.log role: 'exocom', text: "#{bold message.sender}  --[ #{bold message.name} ]->  #{bold receivers.join ' and '}"
      indent = ' ' * (message.sender.length + 2)
      if message.payload?
        for line in util.inspect(message.payload, show-hidden: false, depth: null).split '\n'
          logger.log role: 'exocom', text: "#{indent}#{dim line}", trim: no
      else
        logger.log role: 'exocom', text: "#{indent}#{dim '(no payload)'}", trim: no
    ..start-exocom!
    ..start-services!

  process.on 'SIGINT', ~>
    app-runner.shutdown close-message: " shutting down ..."

function help
  help-message =
    """
    Usage: #{cyan "exo run"}

    Runs an Exosphere application.
    This command must be run in the root directory of the application.
    """
  console.log help-message
