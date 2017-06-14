require! {
  '../../exosphere-shared': {run-process}
  'chai': {expect}
  'fs-extra' : fs
  'js-yaml' : yaml
  'mkdirp'
  'path'
  'prelude-ls': {each}
  'rimraf'
}


describe 'AppSetup', ->

  describe 'set up a complex app' (...) ->
    @timeout 600000

    before (done) ~>
      @app-name = 'complex-setup-app'
      @app-dir = path.join process.cwd!, 'tmp', @app-name
      @docker-compose-location = path.join @app-dir, 'tmp', 'docker-compose.yml'
      @exocom-name = 'exocom0.22.1'
      @internal-services = ['html-server', 'todo-service']
      @external-services = ['external-service']
      @external-dependencies = ['mongo3.4.0']
      @process = checkout-setup-app @app-name, @app-dir
        ..on 'ended' done

    specify 'should create docker-compose.yml at the expected location', (done) ~>
      fs.stat @docker-compose-location, (err, stat) ~>
        expect(err).to.be.null
        @docker-compose = yaml.safe-load fs.read-file-sync(@docker-compose-location, 'utf8')
        done!

    specify 'should list all services and dependencies under \'services\'', (done) ~>
      all-services = [@exocom-name] ++ @internal-services ++ @external-dependencies ++ @external-services
      expect(Object.keys @docker-compose.services).to.have.members(all-services)
      done!

    specify 'should generate an image name for each dependency and external service', (done) ~>
      expect(@docker-compose.services[@exocom-name].image).to.not.be.empty
      @external-dependencies.for-each (dependency) ~> expect(@docker-compose.services[dependency].image).to.not.be.empty
      @external-services.for-each (service) ~> expect(@docker-compose.services[service].image).to.not.be.empty
      done!

    specify 'should generate a container name for each service and dependency', (done) ~>
      expect(@docker-compose.services[@exocom-name].container_name).to.not.be.empty
      @external-dependencies.for-each (dependency) ~> expect(@docker-compose.services[dependency].container_name).to.not.be.empty
      @external-services.for-each (service) ~> expect(@docker-compose.services[service].container_name).to.not.be.empty
      @internal-services.for-each (service) ~> expect(@docker-compose.services[service].container_name).to.not.be.empty
      done!

    specify 'should have the correct build command for each service and dependency', (done) ~>
      @internal-services.for-each (service) ~> expect(@docker-compose.services[service].command).to.eql('echo "does not run"')
      expect(@docker-compose.services[@exocom-name].command).to.eql('bin/exocom')
      done!

    specify 'should have the correct build path for each internal service', (done) ~>
      @internal-services.for-each (service) ~> expect(@docker-compose.services[service].build).to.eql("../#service")
      done!

    specify 'should include \'exocom\' in the dependencies of every service', (done) ~>
      @internal-services ++ @external-services.for-each (service) ~> expect(@docker-compose.services[service].depends_on).to.include(@exocom-name)
      done!

    specify 'should include external dependencies as dependencies', (done) ~>
      expect(@docker-compose.services['todo-service'].depends_on).to.include('mongo3.4.0')
      done!

    specify 'should set up an \'environment\' for exocom with a port and correct service routes', (done) ~>
      exocom-env = @docker-compose.services[@exocom-name].environment
      expect(exocom-env.ROLE).to.eql('exocom')
      expect(exocom-env.PORT).to.eql('$EXOCOM_PORT')
      expect(exocom-env.SERVICE_ROUTES).to.eql('[{"role":"html-server","receives":["todo.created"],"sends":["todo.create"]},{"role":"todo-service","receives":["todo.create"],"sends":["todo.created"]},{"role":"external-service","receives":["users.listed","users.created"],"sends":["users.list","users.create"]}]')
      done!

    specify 'should set up an \'environment\' for internal service with exocom as host', (done) ~>
      @internal-services.for-each (service) ~>
        env = @docker-compose.services[service].environment
        expect(env.ROLE).to.eql(service)
        expect(env.EXOCOM_HOST).to.eql(@exocom-name)
        expect(env.EXOCOM_PORT).to.eql('$EXOCOM_PORT')
      done!

    specify 'should generate a volume path for an external dependency that mounts a volume', (done) ~>
      expect(@docker-compose.services['mongo3.4.0'].volumes).to.not.be.empty
      done!

    specify 'should have the specified image and container names for the external service', (done) ~>
      @service-name = 'external-service'
      @image-name = 'originate/test-web-server'
      expect(@docker-compose.services[@service-name].image).to.eql(@image-name)
      expect(@docker-compose.services[@service-name].container_name).to.eql(@service-name)
      done!


checkout-setup-app = (app-name, app-dir) ->
  rimraf.sync app-dir
  mkdirp.sync app-dir
  fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name),
               app-dir
  run-process path.join(process.cwd!, 'bin', 'exo-setup'), app-dir
