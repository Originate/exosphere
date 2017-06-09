require! {
  'chai' : {expect}
  'cucumber': {defineSupportCode}

  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'path'
}


defineSupportCode ({Then}) ->

  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.readFile path.join(@app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  Then /^my workspace contains the file "([^"]*)" with content:$/, (filename, expected-content, done) ->
    fs.readFile path.join(@app-dir, filename), N (actual-content) ->
      jsdiff-console actual-content.toString!trim!, expected-content.trim!, done
