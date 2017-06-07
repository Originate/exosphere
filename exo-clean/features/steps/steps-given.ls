require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args, DockerHelper}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'wait' : {wait}
}


module.exports = ->

  @Given /^my machine has both dangling and non-dangling Docker images and volumes$/ timeout: 600_000, (done) ->
    @app-name = 'external-dependency'
    @service-name = 'mongo'
    @image-name = 'mongo:3.4.0'
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
        @add-file @app-dir, @service-name, 'test.txt', ~>
          @setup-app @app-dir, ~>
            DockerHelper.start-container {Image: @image-name}, (err) ->
              DockerHelper.remove-container {Image: @image-name}, done
