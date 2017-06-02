require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'request'
  'fs'
}


module.exports = ->

  @When /^running "([^"]*)" in this application's directory$/ timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)" in the terminal$/ timeout: 6_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^trying to run the "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      command = \exo-run
      if process.platform is 'win32' then command += '.cmd'
      @run-app {command}, done


  @When /^the web service broadcasts a "([^"]*)" message$/ (message, done) ->
    request 'http://localhost:4000', done


  @When /^waiting until I see "([^"]*)" in the terminal$/ timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @When /^adding a file to the "([^"]*)" service$/ (service-name) ->
    app-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'application.yml'), 'utf8')
    service-config = app-config.services[\public][service-name] or app-config.services[\private][service-name]
    fs.write-file-sync path.join(@app-dir, service-config.location, 'test.txt'), 'test'
