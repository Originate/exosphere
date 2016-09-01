require! {
  'async'
  'chai' : {expect}
  'dim-console'
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'request'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-name, done


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^the web service broadcasts a "([^"]*)" message$/, (message, done) ->
    request 'http://localhost:4000', done


  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @Then /^ExoCom uses this routing:$/, (table, done) ->
    expected-routes = {}
    for row in table.hashes!
      expected-routes[row.COMMAND] or= {}
      for receiver in row.RECEIVERS.split(', ')
        (expected-routes[row.COMMAND].receivers or= []).push host: 'localhost', name: receiver
    request "http://localhost:8001/config.json", N (response, body) ->
      expect(response.status-code).to.equal 200
      actual-routes = JSON.parse(body).routes
      for _, data of actual-routes
        for receiver in data.receivers
          expect(receiver.port).to.be.at.least 3000
          delete receiver.port
          delete receiver.internal-namespace
      jsdiff-console actual-routes, expected-routes, done


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^my machine is running ExoCom$/, (done) ->
    @process.wait 'exocom  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done


  @Then /^the "([^"]*)" service receives a "([^"]*)" message$/, (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done


  @Then /^the "([^"]*)" service replies with a "([^"]*)" message$/, (arg1, arg2, done) ->
    done!

