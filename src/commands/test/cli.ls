require! {
  './app-tester' : AppTester
  'chalk' : {cyan, green, red}
  'fs'
  'js-yaml' : yaml
  '../../logger' : Logger
  '../../../package.json' : {version}
  'path'
}

app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
logger = new Logger Object.keys(app-config.services)
logger.log name: 'exo-test', text: "Testing application '#{app-config.name}'"
app-tester = new AppTester app-config
  ..on 'output', (data) -> logger.log data
  ..on 'all-tests-passed', -> logger.log name: 'exo-test', text: 'All tests passed'
  ..on 'all-tests-failed', -> logger.log name: 'exo-test', text: 'Tests failed'
  ..on 'service-tests-passed', (name) -> logger.log name: 'exo-test', text: "#{name} service works"
  ..on 'service-tests-failed', (name) -> logger.log name: 'exo-test', text: "#{name} service is broken"
  ..on 'service-tests-skipped', (name) -> logger.log name: 'exo-test', text: "#{name} service has no tests, skipping"
  ..start-testing!
