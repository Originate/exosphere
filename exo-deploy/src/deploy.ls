require! {
  './app-deployer' : AppDeployer
  'fs'
  'js-yaml' : yaml
}

app-config = yaml.safe-load fs.read-file-sync('/var/app/application.yml', 'utf8')
new AppDeployer app-config
  ..start!
