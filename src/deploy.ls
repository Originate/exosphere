require! {
  './app-deployer' : AppDeployer
  'exosphere-shared' : {Logger}
  'require-yaml'
}

app-config = require '/var/app/application.yml'
logger = new Logger
new AppDeployer app-config, logger
  ..start!
