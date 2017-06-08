require! {
  'cucumber': {defineSupportCode}
  'chai' : {expect}
}


defineSupportCode ({then}) ->

  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text
