require! {
  './docker/docker' : Docker
  'chalk' : {cyan, green}
  'path'
  'require-yaml'
}

app-config = require path.join(process.cwd!, 'application.yml')
console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
new Docker
  ..start!
