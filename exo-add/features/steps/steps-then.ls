require! {
  'chai' : {expect}
  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'path'
}


module.exports = ->

  @Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.read-file path.join(@app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  @Then /^my application contains the file "([^"]*)"$/, (file-path) ->
    expect(fs.exists-sync path.join(@app-dir, file-path)).to.be.true


  @Then /^my application contains the file "([^"]*)" containing the text:$/, (file-path, expected-fragment, done) ->
    fs.read-file path.join(@app-dir, file-path), N (actual-content) ->
      expect(actual-content.to-string!).to.contain expected-fragment.trim!
      done!


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text

  @Then /^I see:$/ (expected-text) ->
    expect(@process.full-output!).to.contain expected-text

  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @Then /^it exits with code (\d+)$/ (+expected-exit-code) ->
    @process.on 'ended', ~>
      expect(@process.exit-code).to.equal expected-exit-code


  @Then /^I see the error "([^"]*)"$/, (expected-text, done) ->
    @process.wait expected-text, done
