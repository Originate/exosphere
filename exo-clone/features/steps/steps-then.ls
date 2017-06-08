require! {
  'cucumber': {defineSupportCode} 
  'chai' : {expect}
  'fs-extra' : fs
  'path'
  'prelude-ls' : {difference, unique}
  'strip-ansi'
}


defineSupportCode ({then}) ->

  Then /^it creates the files:$/ (table) ->
    for row in table.raw!
      fs.access-sync path.join process.cwd!, 'tmp', row[0]
      @existing-folders.push(row[0].split(path.sep)[0])


  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(strip-ansi @process.full-output!).to.contain expected-text


  Then /^I get the error "([^"]*)"$/ (expected-text) ->
    expect(strip-ansi @process.full-output!).to.contain expected-text

  Then /^no new files or folders have been created$/ (done) ->
    files-list = fs.readdir-sync path.join process.cwd!, 'tmp' |> unique
    new-files = difference files-list, @existing-folders
    if new-files.length > 0
      done new Error "New files or folders were created: #{new-files}"
    done!
