require! {
  './app-deployer' : AppDeployer
  'exosphere-shared' : {Logger}
  'chalk' : {cyan, green}
  'path'
  'require-yaml'
}

app-config = require path.join(process.cwd!, 'application.yml')
console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger
new AppDeployer app-config, logger
  ..start!
