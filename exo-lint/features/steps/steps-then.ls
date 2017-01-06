require! {
  'chai' : {expect}
}


module.exports = ->

  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text
