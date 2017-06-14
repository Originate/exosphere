require! {
  '../../exosphere-shared': {run-process}
  'chai': {expect}
  'fs-extra' : fs
  'js-yaml' : yaml
  'mkdirp'
  'path'
  'rimraf'
}


describe 'AppSetup', ->

  describe 'app with external dependencies' (...) ->
    @timeout 600000

    before (done) ~>
      @app-name = 'external-dependency'
      @app-dir = path.join process.cwd!, 'test', @app-name
      @docker-compose-location = path.join @app-dir, 'tmp', 'docker-compose.yml'
      @exocom-name = 'exocom0.22.1'
      @service-name = 'mongo'
      @mongo-name = 'mongo3.4.0' # external dependency
      @process = checkout-setup-app @app-name, @app-dir
        ..on 'ended' done

    specify 'should create docker-compose.yml at the expected location', (done) ~>
      fs.exists @docker-compose-location, (exists) ~>
        expect(exists).to.be.true
        @docker-compose = yaml.safe-load fs.read-file-sync(@docker-compose-location, 'utf8')
        done!

    specify 'should list all services and dependencies under \'services\'', (done) ~>
      expect(Object.keys @docker-compose.services).to.eql([@exocom-name, @service-name, @mongo-name])
      done!

    specify 'should generate an image name for each dependency', (done) ~>
      expect(@docker-compose.services[@exocom-name].image).to.not.be.empty
      expect(@docker-compose.services[@mongo-name].image).to.not.be.empty
      done!

    specify 'should generate a container name for each service/dependency equal to the service/dependency name', (done) ~>
      expect(@docker-compose.services[@service-name].container_name).to.eql(@service-name)
      expect(@docker-compose.services[@mongo-name].container_name).to.eql(@mongo-name)
      expect(@docker-compose.services[@exocom-name].container_name).to.eql(@exocom-name)
      done!

    specify 'should have the correct command for each service/dependency that can be built', (done) ~>
      expect(@docker-compose.services[@service-name].command).to.eql('node_modules/exoservice/bin/exo-js')
      expect(@docker-compose.services[@exocom-name].command).to.eql('bin/exocom')
      done!

    specify 'should have correct build path for each service', (done) ~>
      expect(@docker-compose.services[@service-name].build).to.eql('../mongo')
      done!

    specify 'should include \'exocom\' in the dependencies of every service', (done) ~>
      expect(@docker-compose.services[@service-name].depends_on).to.include(@exocom-name)
      done!

    specify 'should include external dependencies', (done) ~>
      expect(@docker-compose.services[@service-name].depends_on).to.include(@mongo-name)
      done!

    specify 'should set up an \'environment\' for exocom with a port and correct service routes', (done) ~>
      exocom-env = @docker-compose.services[@exocom-name].environment
      expect(exocom-env.ROLE).to.eql('exocom')
      expect(exocom-env.PORT).to.eql('$EXOCOM_PORT')
      expect(exocom-env.SERVICE_ROUTES).to.eql('[{"role":"mongo","receives":null,"sends":null}]')
      done!

    specify 'should set up an \'environment\' for each service with exocom as host', (done) ~>
      mongo-env = @docker-compose.services[@service-name].environment
      expect(mongo-env.ROLE).to.eql(@service-name)
      expect(mongo-env.EXOCOM_HOST).to.eql(@exocom-name)
      expect(mongo-env.EXOCOM_PORT).to.eql('$EXOCOM_PORT')
      done!

    specify 'should generate a volume path for an external dependency that mounts a volume', (done) ~>
      expect(@docker-compose.services[@mongo-name].volumes).to.not.be.empty
      done!

  # indirectly tests @cat-file used in exosphere-shared/src/compile-service-routes.ls
  describe 'app with external docker images' (...) ->
    @timeout 600000

    before (done) ~>
      @app-name = 'app-with-external-docker-images'
      @app-dir = path.join process.cwd!, 'test', @app-name
      @docker-compose-location = path.join @app-dir, 'tmp', 'docker-compose.yml'
      @exocom-name = 'exocom0.22.1'
      @service-name = 'external-service'
      @image-name = 'originate/test-web-server'
      @process = checkout-setup-app @app-name, @app-dir
        ..on 'ended', ~>
          @docker-compose = yaml.safe-load fs.read-file-sync(@docker-compose-location, 'utf8')
          done!

    it 'should include external docker image under \'services\'', (done) ~>
      expect(Object.keys @docker-compose.services).to.include(@service-name)
      done!

    it 'should have the specified image and container names for the external docker image', (done) ~>
      expect(@docker-compose.services[@service-name].image).to.eql(@image-name)
      expect(@docker-compose.services[@service-name].container_name).to.eql(@service-name)
      done!

    it 'should include \'exocom\' in the dependencies of external docker image service', (done) ~>
      expect(@docker-compose.services[@service-name].depends_on).to.include(@exocom-name)
      done!

    it 'should include external docker image\'s service route in exocom\'s \'environment\'', (done) ~>
      exocom-env = @docker-compose.services[@exocom-name].environment
      expect(exocom-env.SERVICE_ROUTES).to.eql('[{"role":"external-service","receives":["users.listed","users.created"],"sends":["users.list","users.create"]}]')
      done!


checkout-setup-app = (app-name, app-dir) ->
  rimraf.sync app-dir
  mkdirp.sync app-dir
  fs.copy-sync path.join(process.cwd!, '..' 'exosphere-shared' 'example-apps', app-name),
               app-dir
  run-process path.join(process.cwd!, 'bin', 'exo-setup'), app-dir
