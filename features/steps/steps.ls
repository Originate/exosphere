require! {
  'async'
  'chai' : {expect}
  '../support/dim-console'
  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'request'
  'tmp'
}


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name


  @Given /^a set\-up "([^"]*)" application$/, timeout: 60*1000, (@app-name, done) ->
    @checkout-app @app-name
    @setup-app @app-name, done



  @When /^installing it$/, timeout: 60_000, (done) ->
    @setup-app @app-name, done


  @When /^starting it$/, timeout: 10_000, (done) ->
    @start-app @app-name, done


  @When /^starting the "([^"]*)" application$/, (@app-name, done) ->
    @start-app @app-name, done



  @Then /^ExoComm uses this routing:$/, (table, done) ->
    expected-routes = {}
    for row in table.hashes!
      expected-routes[row.COMMAND] or= {}
      for receiver in row.RECEIVERS.split(', ')
        (expected-routes[row.COMMAND].receivers or= []).push host: 'localhost', name: receiver
    request "http://localhost:5000/config.json", N (response, body) ->
      expect(response.status-code).to.equal 200
      actual-routes = JSON.parse(body).routes
      for _, data of actual-routes
        for receiver in data.receivers
          expect(receiver.port).to.be.at.least 3000
          delete receiver.port
      jsdiff-console actual-routes, expected-routes, done


  @Then /^it creates the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@app-dir, row.SERVICE, row.FOLDER), fs.F_OK


  @Then /^my machine is running ExoComm$/, (done) ->
    @process.wait 'exocomm  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done
