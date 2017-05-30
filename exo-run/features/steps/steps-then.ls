require! {
  'async'
  'chai' : {expect}
  'child_process'
  '../../../exosphere-shared' : {compile-service-routes}
  'fs'
  'jsonic'
  'jsdiff-console'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'prelude-ls' : {any, last}
  'os'
  'path'
  'request'
  'wait' : {wait}
}


module.exports = ->

  @Then /^ExoCom uses this routing:$/ (table) ->
    expected-routes = []
    for row in table.hashes!
      service-routes = {}
      service-routes.role = row.ROLE
      for message in row.RECEIVES.split(', ')
        (service-routes.receives or= []).push message
      for message in row.SENDS.split(', ')
        (service-routes.sends or= []).push message
      if row.NAMESPACE
        service-routes.namespace = row.NAMESPACE
      expected-routes.push service-routes
    docker-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'docker-compose.yml'))
    actual-routes = jsonic docker-config.services.exocom.environment.SERVICE_ROUTES
    jsdiff-console actual-routes, expected-routes


  @Then /^it has printed "([^"]*)" in the terminal$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^I see:$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^it prints "([^"]*)" in the terminal$/ timeout: 60_000, (expected-text, done) ->
    @process.wait expected-text, done


  @Then /^the "([^"]*)" service restarts$/ (service, done) ->
    @process.wait "Restarting service '#{service}'", done


  @Then /^my machine is running ExoCom$/ timeout: 10_000, (done) ->
    @process.wait /ExoCom WebSocket listener online at port/, done


  @Then /^my machine is running the services:$/ (table) ->
    for row in table.hashes!
      expect(child_process.exec-sync('docker ps --format {{.Names}}/{{.Status}}') |> (.to-string!) |> (.split os.EOL) |> any (.includes "#{row.NAME}/Up")).to.be.true


  @Then /^the "([^"]*)" service receives a "([^"]*)" message$/ (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done

