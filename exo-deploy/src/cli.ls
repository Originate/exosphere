require! {
  'chalk' : {cyan, green}
  './docker-runner' : Docker
  '../../exosphere-shared' : {Logger}
  'fs'
  'inquirer'
  'js-yaml' : yaml
  'path'
}

module.exports = ->

  if process.argv[2] is \help
    return help!

  command-flag = process.argv[2]
  app-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, 'application.yml'), 'utf8')
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
        | err => return logger.log role: 'exo-deploy', text: err.message
        ..start!

function help
  help-text = """
  Usage: #{cyan "exo deploy"}

  """
  console.log help-text
