# NOTE: this file is run inside a docker container
require! {
  './app-deployer' : AppDeployer
  'require-yaml'
}

command-flag = process.argv[2]
app-config = require '/var/app/application.yml'
deployer = new AppDeployer app-config, command-flag
if command-flag is '--nuke' then
  deployer.nuke!
else
  deployer.deploy!
