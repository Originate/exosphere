require! {
  'async'
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



module.exports = ->

  @When /^I start the "([^"]*)" application$/, (app-name, done) ->
    @process = new ObservableProcess(path.join('..', '..', 'bin', 'exo-run'),
                                     cwd: path.join(process.cwd!, 'example-apps', app-name),
                                     verbose: yes,
                                     console: dim-console)
      ..wait 'all systems online', done


  @Then /^my machine is running ExoComm$/, (done) ->
    @process.wait 'exocomm  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done
