# NOTE: this file is run inside a docker container
require! {
  './app-deployer' : AppDeployer
  'fs'
  'js-yaml' : yaml
}

command-flag = process.argv[2]
app-config = yaml.safe-load fs.read-file-sync('/var/app/application.yml', 'utf8')
deployer = new AppDeployer app-config
if command-flag is '--nuke' then
  deployer.teardown nuke: yes, (err) ->
    | err => process.stdout.write err.message
else if command-flag is '--teardown' then
  deployer.teardown nuke: no, (err) ->
    | err => process.stdout.write err.message
else
  deployer.deploy (err) ->
    | err => process.stdout.write err.message
