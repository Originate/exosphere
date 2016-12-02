require! {
  './app-deployer' : AppDeployer
  'require-yaml'
}

app-config = require '/var/app/application.yml'
new AppDeployer app-config
  ..start!
