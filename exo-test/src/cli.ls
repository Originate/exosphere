require! {
  './app-tester' : AppTester
  'chalk' : {cyan, green, red}
  'fs'
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'path'
  './service-tester' : ServiceTester
}

module.exports = ->

  switch
    | cwd-is-service! => test-service!
    | cwd-is-app! => test-app!
    | otherwise => logger = new Logger!.log name: 'exo-test', text: "Tests do not exist. Not in service or application directory."

function cwd-is-service
  try
    fs.stat-sync 'service.yml'
  catch
    false

function cwd-is-app
  try
    fs.stat-sync 'application.yml'
  catch
    false

function test-service
  service-name = path.basename process.cwd!
  logger = new Logger [service-name]
    ..log name: 'exo-test', text: "Testing service '#{service-name}'"
  new ServiceTester {name: service-name, config: {root: process.cwd!}, logger}
    ..start ~>
      ..remove-dependencies!

function test-app
  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  logger = new Logger Object.keys(app-config.services)
    ..log name: 'exo-test', text: "Testing application '#{app-config.name}'"
  app-tester = new AppTester {app-config, logger}
    ..start-testing!
