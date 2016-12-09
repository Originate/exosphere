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
  new ServiceTester service-name, root: process.cwd!
    ..on 'output', (data) ~> logger.log data
    ..on 'service-tests-passed', (name) -> logger.log name: 'exo-test', text: "#{name} works"
    ..on 'service-tests-failed', (name) -> logger.log name: 'exo-test', text: "#{name} is broken"
    ..on 'service-tests-skipped', (name) -> logger.log name: 'exo-test', text: "#{name} has no tests, skipping"
    ..start!

function test-app
  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  logger = new Logger Object.keys(app-config.services)
    ..log name: 'exo-test', text: "Testing application '#{app-config.name}'"
  app-tester = new AppTester app-config
    ..on 'output', (data) -> logger.log data
    ..on 'all-tests-passed', -> logger.log name: 'exo-test', text: 'All tests passed'
    ..on 'all-tests-failed', -> logger.log name: 'exo-test', text: 'Tests failed'; process.exit 1
    ..on 'service-tests-passed', (name) -> logger.log name: 'exo-test', text: "#{name} works"
    ..on 'service-tests-failed', (name) -> logger.log name: 'exo-test', text: "#{name} is broken"
    ..on 'service-tests-skipped', (name) -> logger.log name: 'exo-test', text: "#{name} has no tests, skipping"
    ..start-testing!
