# NOTE: this file is run inside a docker container
require! {
  './app-deployer' : AppDeployer
  'fs'
  'js-yaml' : yaml
}

command-flag = process.argv[2]
app-config = require '/var/app/application.yml'
deployer = new AppDeployer app-config, command-flag
if command-flag is '--nuke' then
  deployer.teardown nuke: yes, (err) ->
    | err => process.stdout.write err.message
else if command-flag is '--teardown' then
  deployer.teardown nuke: no, (err) ->
    | err => process.stdout.write err.message
else
  deployer.deploy (err) ->
    | err => process.stdout.write err.message
