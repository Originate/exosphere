require! {
  'chalk' : {cyan, green}
  './docker/docker' : Docker
  'exosphere-shared' : {Logger}
  'inquirer'
  'path'
  'require-yaml'
}

command-flag = process.argv[2]
app-config = require path.join(process.cwd!, 'application.yml')
logger = new Logger
docker = new Docker app-config, logger
if command-flag is '--nuke'
  console.log "You are about to nuke #{green app-config.name} #{cyan app-config.version}\n"

  question =
    type: 'list'
    name: 'continue'
    message: 'Are you sure?'
    choices: ['yes', 'no']
  inquirer.prompt([question]).then (answer) ->
    if answer.continue == 'no'
      console.log '\nAborting ...\n'
      process.exit!
    else docker.start command-flag
else
  console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
  docker
    ..dockerhub-push (err) ~>
      | err => return logger.log name: 'exo-deploy', text: err.message
      ..start!
