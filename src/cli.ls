require! {
  './app-deployer' : AppDeployer
  'chalk' : {bold, cyan, dim, green, red}
  'fs'
  'js-yaml' : yaml
  'exosphere-shared' : {Logger}
}

app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
console.log "Deploying #{green app-config.name} #{cyan app-config.version}\n"
logger = new Logger
app-deployer = new AppDeployer app-config
  ..deploy!
