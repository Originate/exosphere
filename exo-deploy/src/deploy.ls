require! {
  './app-deployer' : AppDeployer
  'require-yaml'
}

command-flag = process.argv[2]
console.log "command: #{command-flag}"
app-config = require '/var/app/application.yml'
deployer = new AppDeployer app-config, command-flag
if command-flag is '--teardown' then
  deployer.teardown nuke: no
else if command-flag is '--nuke' then
  deployer.teardown nuke: yes
else
  deployer.deploy!
