require! {
  'chai' : {expect}
  '../support/dim-console'
  'observable-process' : ObservableProcess
  'path'
  'request'
}


check-exocomm-port = (port) ->
  request-data =
    url: "http://localhost:#{port}/send/foo"
    method: 'POST'
    body:
      request-id: '123'
    json: yes
  request request-data, (err, response) ->
    | err  =>  throw new Error "Expected app '#{app-name}' to have ExoRelay port #{port}"



check-exorelay-port = (app-name, port) ->
  request-data =
    url: "http://localhost:#{port}/run/test"
    method: 'POST'
    body:
      request-id: '123'
    json: yes
  request request-data, (err, response) ->
    | err  =>  throw new Error "Expected app '#{app-name}' to have ExoRelay port #{port}"


check-public-port = (port) ->
  request "http://localhost:#{port}/run/test", (err, response) ->
    console.log err
    expect(err).to.be.null


module.exports = ->

  @When /^I start the "([^"]*)" application$/, (app-name, done) ->
    @process = new ObservableProcess(path.join('..', '..', 'bin', 'exo-run'),
                                     cwd: path.join(process.cwd!, 'example-apps', app-name),
                                     verbose: yes,
                                     console: dim-console)
      ..wait 'all systems online', done


  @Then /^my machine is running ExoComm at port (\d+)$/, (port) ->
    check-exocomm-port port


  @Then /^my machine is running the services:$/, (table, done) ->
    for app in table.hashes!
      check-exorelay-port app['NAME'], app['EXORELAY-PORT']
      check-public-port app['PUBLIC PORT']
