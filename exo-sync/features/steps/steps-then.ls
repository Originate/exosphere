require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}

  'fs'
  'path'
}


defineSupportCode ({Then}) ->

  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^my application contains the newly committed file$/, ->
    fs.stat-sync path.join(@app-dir, 'web-service', 'new_file')
