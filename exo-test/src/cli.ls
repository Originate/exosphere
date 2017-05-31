require! {
  './app-tester' : AppTester
  'chalk' : {cyan, green, red}
  'fs'
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'path'
  'prelude-ls' : {flatten}
  './service-tester' : ServiceTester
}

module.exports = ->

  if process.argv[2] is "help"
    return help!

  switch
    | cwd-is-service! => test-service!
    | cwd-is-app! => test-app!
    | otherwise => logger = new Logger!.log role: 'exo-test', text: "Tests do not exist. Not in service or application directory."

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
  service-role = path.basename process.cwd!
  logger = new Logger [service-role]
    ..log role: 'exo-test', text: "Testing service '#{service-role}'"
  new ServiceTester {role: service-role, service-location: process.cwd!, logger}
    ..start!

function test-app
  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  logger = new Logger flatten [Object.keys(app-config.services[protection-level]) for protection-level of app-config.services]
    ..log role: 'exo-test', text: "Testing application '#{app-config.name}'"
  app-tester = new AppTester {app-config, logger}
    ..start-testing!

function help
  help-message =
    """
    \nUsage: #{cyan 'exo test'}

    Runs feature tests for a service or application.
    This command must be called in either a service directory or the root directory of the application.
    """
  console.log help-message
