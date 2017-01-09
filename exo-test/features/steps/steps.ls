require! {
  'chai' : {expect}
  'dim-console'
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-name, done


  @Given /^I am in the "([^"]*)" directory$/ (service-dir, done) ->
    @checkout-app service-dir.split(path.sep)[0]
    @current-dir = path.join process.cwd!, 'tmp', service-dir
    done!


  @Given /^I am in the "([^"]*)" created directory$/ (dir, done) ->
    new-dir = path.join process.cwd!, dir
    try
      fs.mkdir-sync new-dir
    @current-dir = new-dir
    done!


  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @Then /^it doesn't run any tests$/ (done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    expect(@process.full-output!).to.not.include "Testing service"
    @process.wait "exo-test  Tests do not exist. Not in service or application directory.", done


  @Then /^it only runs tests for "([^"]*)"$/ (service-role, done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    @process.wait "exo-test  Testing service '#{service-role}'", done


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


