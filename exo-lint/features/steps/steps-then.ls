require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}

}


defineSupportCode ({Then}) ->

  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text
