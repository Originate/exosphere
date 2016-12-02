require! {
  'chalk' : {cyan, green}
  './docker/docker' : Docker
  'exosphere-shared' : {Logger}
  'path'
  'require-yaml'
}

app-config = require path.join(process.cwd!, 'application.yml')
console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger
new Docker app-config, logger
  ..dockerhub-push (err) ~>
    | err => return logger.log name: 'exo-deploy', text: err.message
    ..start!
