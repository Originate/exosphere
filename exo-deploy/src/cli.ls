require! {
  'chalk' : {cyan, green}
  './docker/docker' : Docker
  '../../exosphere-shared' : {Logger}
  'fs'
  'js-yaml' : yaml
  'path'
}

deploy = ->

  app-config =  yaml.safe-load fs.read-file-sync(path.join(process.cwd!, 'application.yml'), 'utf8')
  console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
  logger = new Logger
  new Docker app-config, logger
    ..dockerhub-push (err) ~>
      | err => return logger.log name: 'exo-deploy', text: err.message
      ..start!



module.exports = deploy
