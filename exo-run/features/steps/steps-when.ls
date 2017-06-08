require! {
  'chai' : {expect}
  '../../../exosphere-shared' : {run-process}
  'js-yaml' : yaml
  'fs-extra' : fs
  'path'
  'request'
  'fs'
  'cucumber': {defineSupportCode}
}


defineSupportCode ({When}) ->


  When /^running "([^"]*)" in this application's directory$/ timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @app-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.be. 0
        done!


  When /^running "([^"]*)" in the terminal$/ timeout: 6_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @app-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.be. 0
        done!


  When /^trying to run the "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-and-run-app {}, done


  When /^the web service broadcasts a "([^"]*)" message$/ (message, done) ->
    request 'http://localhost:4000', done


  When /^waiting until I see "([^"]*)" in the terminal$/ timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  When /^adding a file to the "([^"]*)" service$/ (service-name) ->
    app-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'application.yml'), 'utf8')
    service-config = app-config.services[\public][service-name] or app-config.services[\private][service-name]
    fs.write-file-sync path.join(@app-dir, service-config.location, 'test.txt'), 'test'
