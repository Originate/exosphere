require! {
  'chai' : {expect}
  'child_process'
  'dim-console'
  '../../../exosphere-shared' : {call-args, DockerHelper}
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    @current-dir = path.join process.cwd!, 'tmp', @app-name


  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', -> done!


  @Then /^it has created the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@current-dir, row.SERVICE, row.FOLDER)


  @Then /^it has created the files:$/, (table) ->
    for row in table.raw!
      fs.access-sync path.join(@current-dir, row[0])


  @Then /^it has acquired the Docker images:$/ (table) ->
    docker-images = DockerHelper.get-docker-images! |> (.to-string!)
    for row in table.raw!
      expect(docker-images).to.include row[0]


  @Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code
