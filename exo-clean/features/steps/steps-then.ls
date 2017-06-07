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


  Then /^it has non-dangling images$/ (done) ->
    DockerHelper.list-images (err, all-images) ->
      DockerHelper.get-dangling-images (err, dangling-images) ->
        expect(all-images.length).to.be.greater-than dangling-images.length 
        done!


  Then /^it does not have dangling images/ (done) ->
    DockerHelper.get-dangling-images (err, images) ->
      expect(images.length).to.equal 0
      done!

  Then /^it does not have dangling volumes$/ (done) ->
    DockerHelper.get-dangling-volumes (err, volumes) ->
      expect(volumes.length).to.equal 0
      done!
