require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}

}


defineSupportCode ({Then}) ->

  Then /^it doesn't run any tests$/ (done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    expect(@process.full-output!).to.not.include "Testing service"
    @process.wait "exo-test  Tests do not exist. Not in service or application directory.", done


  Then /^it only runs tests for "([^"]*)"$/ (service-name, done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    @process.wait "exo-test  Testing service '#{service-name}'", done


  Then /^it prints "([^"]*)" in the terminal$/, (expected-text, done) ->
    @process.wait expected-text, done
