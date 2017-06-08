require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}

  '../../../exosphere-shared' : {DockerHelper, compile-service-routes}
  'fs'
  'jsdiff-console'
  'js-yaml' : yaml
  'path'
}


defineSupportCode ({Then}) ->

  Then /^it has created the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@current-dir, row.SERVICE, row.FOLDER)


  Then /^it has created the files:$/, (table) ->
    for row in table.raw!
      fs.access-sync path.join(@current-dir, row[0])


  Then /^it has acquired the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      for row in table.raw!
        expect(docker-images).to.include row[0]
      done!


  Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code


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
    docker-config = yaml.safe-load fs.read-file-sync(path.join(@current-dir, 'tmp', 'docker-compose.yml'))
    actual-routes = JSON.parse docker-config.services['exocom0.21.8'].environment.SERVICE_ROUTES
    jsdiff-console actual-routes, expected-routes
