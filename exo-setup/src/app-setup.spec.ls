require! {
  './app-setup': AppSetup
  'chai': {expect}
  'fs'
  'js-yaml' : yaml
  '../../exosphere-shared': {Logger}
  'path'
  'prelude-ls' : {flatten}
}

describe 'AppSetup', ->

  # app-path = path.join '../exosphere-shared', 'example-apps', 'simple'
  # app-config = yaml.safe-load fs.read-file-sync("#app-path/application.yml", 'utf8')
  # logger = new Logger flatten [Object.keys(app-config.services[type]) if app-config.services[type] for type of app-config.services]
  # app-setup = new AppSetup app-config: app-config, logger: logger
  # app-setup._get-dependencies-docker-config!
  # console.log app-setup.docker-compose-config

  describe '_get-dependencies-docker-config', ->

    specify 'should parse non-service log message correctly' ->
      expect(0).to.eql(0)

  describe '_get-service-docker-config', ->
    specify 'should parse non-service log message correctly' ->
      expect(0).to.eql(0)

  describe '_render-docker-compose', ->
    specify 'should parse non-service log message correctly' ->
      expect(0).to.eql(0)
