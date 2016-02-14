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
    @app-dir = tmp.dir-sync!
    fs.copy-sync path.join(process.cwd!, 'example-apps', @app-name), @app-dir.name



  @When /^starting the "([^"]*)" application$/, (app-name, done) ->
    @process = new ObservableProcess(path.join('..', '..', 'bin', 'exo-run'),
                                     cwd: path.join(process.cwd!, 'example-apps', app-name),
                                     verbose: yes,
                                     console: dim-console)
      ..wait "application ready", done



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


  @Then /^my machine is running ExoComm$/, (done) ->
    @process.wait 'exocomm  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done
