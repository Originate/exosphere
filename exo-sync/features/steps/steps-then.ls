require! {
  'cucumber': {defineSupportCode}
  'chai' : {expect}
  'fs'
  'path'
}


defineSupportCode ({then}) ->

  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^my application contains the newly committed file$/, ->
    fs.stat-sync path.join(@app-dir, 'web-service', 'new_file')
