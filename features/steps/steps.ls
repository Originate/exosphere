require! {
  'async'
  'chai' : {expect}
  'child_process'
  'dim-console'
  'exosphere-shared' : {call-args}
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'prelude-ls' : {last}
  'request'
  'fs'
  'wait' : {wait}
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^a set\-up "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name
    done!


  @When /^starting "([^"]*)" in this application's directory$/ (command) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^running "([^"]*)" in this application's directory$/ timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^the web service broadcasts a "([^"]*)" message$/ (message, done) ->
    request 'http://localhost:4000', done


  @When /^waiting until I see "([^"]*)" in the terminal$/ timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @Then /^ExoCom uses this routing:$/ timeout: 10_000, (table, done) ->
    expected-routes = {}
    for row in table.hashes!
      expected-routes[row.COMMAND] or= {}
      for receiver in row.RECEIVERS.split(', ')
        (expected-routes[row.COMMAND].receivers or= []).push name: receiver
    exocom-port = child_process.exec-sync('docker port exocom') |> (.to-string!) |> (.split ':') |> last |> (.trim!)
    wait 10, ~> # Wait to ensure services have time to be registered by ExoCom
      request "http://localhost:#{exocom-port}/config.json", N (response, body) ->
        expect(response.status-code).to.equal 200
        actual-routes = JSON.parse(body).routes
        for _, data of actual-routes
          for receiver in data.receivers
            delete receiver.port
            delete receiver.internal-namespace
        jsdiff-console actual-routes, expected-routes, done


  @Then /^it prints "([^"]*)" in the terminal$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^my machine is running ExoCom$/ timeout: 10_000, (done) ->
    #TODO: Use updated 'wait' that allows for regex as opposed to hardcoded version
    @process.wait 'exocom  ExoCom 0.15.1 WebSocket listener online at port', done


  @Then /^my machine is running the services:$/ timeout: 10_000, (table, done) ->
    async.each [row['NAME'] for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name.to-lower-case!}' is running", cb),
               done


  @Then /^the "([^"]*)" service receives a "([^"]*)" message$/ (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done


  @Then /^the "([^"]*)" service replies with a "([^"]*)" message$/ (arg1, arg2, done) ->
    done!
