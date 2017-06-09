require! {
  'async'
  'chai' : {expect}
  'child_process'
  '../../../exosphere-shared' : {DockerHelper, compile-service-routes}
  'fs'
  'jsdiff-console'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'prelude-ls' : {any, last}
  'os'
  'path'
  'request'
  'wait' : {wait}
  'cucumber': {defineSupportCode}
}


defineSupportCode ({Then}) ->

  Then /^ExoCom uses this routing:$/ (table) ->
    expected-routes = []
    for row in table.hashes!
      service-routes = {role: row.ROLE}
      for message in row.RECEIVES.split(', ')
        (service-routes.receives or= []).push message
      for message in row.SENDS.split(', ')
        (service-routes.sends or= []).push message
      if row.NAMESPACE
        service-routes.namespace = row.NAMESPACE
      expected-routes.push service-routes
    docker-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'tmp', 'docker-compose.yml'))
    actual-routes = JSON.parse docker-config.services['exocom0.22.1'].environment.SERVICE_ROUTES
    jsdiff-console actual-routes, expected-routes


  Then /^it has printed "([^"]*)" in the terminal$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^I see:$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it prints "([^"]*)" in the terminal$/ timeout: 60_000, (expected-text, done) ->
    @process.wait expected-text, done


  Then /^the "([^"]*)" service restarts$/ (service, done) ->
    @process.wait "Restarting service '#{service}'", done


  Then /^my machine is running ExoCom$/ timeout: 10_000, (done) ->
    @process.wait /ExoCom WebSocket listener online at port/, done


  Then /^my machine is running the services:$/ (table, done) ->
    DockerHelper.list-running-containers (err, running-containers) ->
      for row in table.hashes!
        expect(running-containers).to.include row.NAME
      done!


  Then /^the "([^"]*)" service receives a "([^"]*)" message$/ (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done
