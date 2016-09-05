require! {
  'chai' : {expect}
  'dim-console'
  'exosphere-shared' : {bash-path}
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    @current-dir = path.join process.cwd!, 'tmp', @app-name


  @Given /^I am in the directory of the "([^"]*)" application$/ (@app-name) ->
    @current-dir = path.join process.cwd!, 'node_modules' 'exosphere-shared' 'example-apps', @app-name


  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(['bash', '-c', bash-path(path.join process.cwd!, 'bin', command)],
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(['bash', '-c', bash-path(path.join process.cwd!, 'bin', command)],
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', -> done!


  @Then /^it has created the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@current-dir, row.SERVICE, row.FOLDER), fs.F_OK


  @Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code

