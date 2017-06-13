require! {
  './app-setup': AppSetup
  '../../exosphere-shared': {Logger, run-process}
  '../features/support/world': World
  'chai': {expect}
  'fs-extra' : fs
  'js-yaml' : yaml
  'mkdirp'
  'path'
  'prelude-ls' : {flatten}
  'rimraf'
}


describe 'AppSetup', ->
  app-name = 'external-dependency'
  app-dir = path.join process.cwd!, 'test', app-name
  docker-compose-location = path.join app-dir, 'tmp', 'docker-compose.yml'
  exocom-name = 'exocom0.22.1'
  service-name = 'mongo'
  mongo-name = 'mongo3.4.0'

  specify 'should create docker-compose.yml at the expected location' ->
    checkout-setup-app app-name, app-dir, ->
      fs.exists docker-compose-location, (exists) ->
        expect(exists).to.be.true
        @docker-compose = yaml.safe-load fs.read-file-sync(docker-compose-location, 'utf8')

  specify 'should list all services and dependencies under \'services\'' ->
    expect(Object.keys docker-compose.services).to.eql([exocom-name, service-name, mongo-name])

  specify 'should generate an image name for each dependency' ->
    expect(docker-compose.services[exocom-name].image).to.not.be.empty
    expect(docker-compose.services[mongo-name].image).to.not.be.empty

  specify 'should generate a container name for each service/dependency equal to the service/dependency name' ->
    expect(docker-compose.services[service-name].container_name).to.eql(service-name)
    expect(docker-compose.services[mongo-name].container_name).to.eql(mongo-name)
    expect(docker-compose.services[exocom-name].container_name).to.eql(exocom-name)

  specify 'should have the correct command for each service/dependency that can be built' ->
    expect(docker-compose.services[service-name].command).to.eql('node_modules/exoservice/bin/exo-js')
    expect(docker-compose.services[exocom-name].command).to.eql('bin/exocom')

  specify 'should have correct build path for each service' ->
    expect(docker-compose.services[service-name].build).to.eql('../mongo')

  specify 'should include \'exocom\' in the dependencies of every service' ->
    expect(docker-compose.services[service-name].depends_on).to.include(exocom-name)

  specify 'should include external dependencies' ->
    expect(docker-compose.services[service-name].depends_on).to.include(mongo-name)

  specify 'should set up an \'environment\' for exocom with a port and correct service routes' ->
    exocom-env = docker-compose.services[exocom-name].environment
    expect(exocom-env.ROLE).to.eql('exocom')
    expect(exocom-env.PORT).to.eql('$EXOCOM_PORT')
    expect(exocom-env.SERVICE_ROUTES).to.eql('[{"role":"mongo","receives":null,"sends":null}]')

  specify 'should set up an \'environment\' for each service with exocom as host' ->
    mongo-env = docker-compose.services[service-name].environment
    expect(mongo-env.ROLE).to.eql(service-name)
    expect(mongo-env.EXOCOM_HOST).to.eql(exocom-name)
    expect(mongo-env.EXOCOM_PORT).to.eql('$EXOCOM_PORT')

  specify 'should generate a volume path for an external dependency that mounts a volume' ->
    expect(docker-compose.services[mongo-name].volumes).to.not.be.empty


checkout-setup-app = (app-name, app-dir, done) ->
  rimraf.sync app-dir
  mkdirp.sync app-dir
  fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name),
               app-dir
  @process = run-process path.join(process.cwd!, 'bin/exo-setup'), app-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.equal 0
        done!
