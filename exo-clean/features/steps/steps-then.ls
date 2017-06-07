require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}
  '../../../exosphere-shared' : {DockerHelper}
}


defineSupportCode ({Then}) ->

  Then /^I see:$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it prints "([^"]*)" in the terminal$/ timeout: 60_000, (expected-text, done) ->
    @process.wait expected-text, done


  Then /^it has the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      for row in table.raw!
        expect(docker-images).to.include row[0]
      done!


  Then /^it does not have the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      for row in table.raw!
        expect(docker-images).to.not.include row[0]
      done!


  Then /^it does not have dangling volumes$/ (done) ->
    DockerHelper.get-dangling-volumes (err, volumes) ->
      expect(volumes.length).to.equal 0
      done!
