require! {
  'chai' : {expect}
  '../../../exosphere-shared' : {DockerHelper, compile-service-routes}
}


module.exports = ->


  @Then /^I see:$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^it prints "([^"]*)" in the terminal$/ timeout: 60_000, (expected-text, done) ->
    @process.wait expected-text, done


  @Then /^the "([^"]*)" service restarts$/ (service, done) ->
    @process.wait "Restarting service '#{service}'", done


  @Then /^it has the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      for row in table.raw!
        expect(docker-images).to.include row[0]
      done!


  @Then /^it does not have the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      for row in table.raw!
        expect(docker-images).to.not.include row[0]
      done!


  @Then /^my machine has a number of dangling and non-dangling Docker images$/ timeout: 60_000, (done) ->
    DockerHelper.list-images (err, docker-images) ->
      @num-docker-images = docker-images.length
      @num-dangling-images = 0
      for let image in docker-images
        if image == '<none>'
          @num-dangling-images += 1
      done!

  @Then /^only dangling Docker images are removed$/ (done) ->
    DockerHelper.list-images (err, docker-images) ->
      expect(docker-images.length).to.equal(@num-docker-images - @num-dangling-images)
      done!
