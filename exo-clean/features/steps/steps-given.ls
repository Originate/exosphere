require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'wait' : {wait}
}


module.exports = ->

  # @Given /^a running "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
  #   @checkout-app @app-name
  #   @app-dir := path.join process.cwd!, 'tmp', @app-name
  #   @setup-app @app-dir, ~>
  #     command = \exo-run
  #     if process.platform is 'win32' then command += '.cmd'
  #     @process = new ObservableProcess(call-args(path.join process.cwd!, '../exo-run/bin', command),
  #                                      cwd: @app-dir,
  #                                      stdout: dim-console.process.stdout
  #                                      stderr: dim-console.process.stderr)
  #     wait 10_000, done


  @Given /^my machine has both dangling and non-dangling Docker images$/ timeout: 600_000, (done) ->
    @app-name = 'simple'
    @service-name = 'web'
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      @add-test-file @app-dir, @service-name, ~>
        @setup-app @app-dir, done

