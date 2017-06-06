require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'wait' : {wait}
}


module.exports = ->

  @Given /^my machine has both dangling and non-dangling Docker images$/ timeout: 600_000, (done) ->
    @app-name = 'simple'
    @service-name = 'web'
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      @add-file @app-dir, @service-name, 'test.txt', ~>
        @setup-app @app-dir, ~>
          @add-file @app-dir, @service-name, 'config.txt', ~>
            @setup-app @app-dir, done
