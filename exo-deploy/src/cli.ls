require! {
  'chalk' : {cyan, green}
  './docker/docker' : Docker
  'exosphere-shared' : {Logger}
  'inquirer'
  'path'
}

module.exports = ->
  command-flag = process.argv[2]
  app-config = require path.join(process.cwd!, 'application.yml')
  logger = new Logger
  docker = new Docker app-config, logger
  if command-flag is '--nuke'
    console.log "You are about to completely remove all parts of #{green app-config.name} #{cyan app-config.version}\n"

    command-flag = process.argv[2]
    app-config = require path.join(process.cwd!, 'application.yml')
    logger = new Logger
    docker = new Docker app-config, logger
    if (command-flag is '--nuke') or (command-flag is '--teardown')
      action = switch command-flag
      | '--nuke'     => 'completely remove all parts of'
      | '--teardown' => 'teardown'
      console.log "You are about to #{action} #{green app-config.name} #{cyan app-config.version}\n"

      question =
        type: 'list'
        name: 'continue'
        message: 'Are you sure?'
        choices: ['yes', 'no']
      inquirer.prompt([question]).then (answer) ->
        if answer.continue == 'no'
          process.exit 2
        else docker.start command-flag
    else
      console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
      docker
        ..dockerhub-push (err) ~>
          | err => return logger.log name: 'exo-deploy', text: err.message
          ..start!
