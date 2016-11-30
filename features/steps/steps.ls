require! {
  'chai' : {expect}
  'dim-console'
  'exosphere-shared' : {call-args}
  'filendir' : {write-file}
  'fs-extra' : fs
  'mkdirp'
  'observable-process' : ObservableProcess
  'path'
  'prelude-ls' : {difference, unique}
  'strip-ansi'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null

module.exports = ->

  @Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @existing-folders = []


  @Given /^source control contains a "([^"]*)" service$/ (service-name, done) ->
    repo-dir = path.join process.cwd!, 'tmp', 'origins', service-name
    mkdirp repo-dir, (err) ~>
      | err  =>  return done err
      fs.write-file-sync path.join(repo-dir, 'service.yml'), ' '
      @make-repo repo-dir
      done!


  @Given /^source control contains a repo "([^"]*)" with a file "([^"]*)" and the content:$/ (app-name, file-name, file-content, done) ->
    @existing-folders.push('origins')
    repo-dir = path.join process.cwd!, 'tmp', 'origins', app-name
    write-file path.join(repo-dir, file-name), file-content, (err) ~>
      @make-repo repo-dir
      done err


  @Given /^source control contains the services "([^"]*)" and "([^"]*)"$/ (service1, service2, done) ->
    @existing-folders.push('origins')
    repo-dirs = [ path.join(process.cwd!, 'tmp', 'origins', service1),
                  path.join(process.cwd!, 'tmp', 'origins', service2) ]
    for dir in repo-dirs
      mkdirp.sync dir
      file-name = path.join dir, 'service.yml'
      fs.write-file-sync file-name, ''
      @make-repo dir
    done!


  @Given /^my workspace already contains the folder "([^"]*)"$/ (folder-name, done) ->
    @existing-folders.push(folder-name)
    temp-file = path.join 'tmp', folder-name, 'temp-file'
    write-file temp-file, ' ', done


  @When /^running "([^"]*)" in the terminal$/, timeout: 10_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, 'tmp')
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', (err, exit-code) ->
        done err


  @When /^trying to run "([^"]*)"$/ (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, 'tmp')
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', (err, exit-code) ~>
        expect(strip-ansi @process.full-output!).to.contain 'Error'
        done!


  @Then /^it creates the files:$/ (table) ->
    for row in table.raw!
      fs.access-sync path.join process.cwd!, 'tmp', row[0]
      @existing-folders.push(row[0].split(path.sep)[0])


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(strip-ansi @process.full-output!).to.contain expected-text


  @Then /^I get the error "([^"]*)"$/ (expected-text) ->
    expect(strip-ansi @process.full-output!).to.contain expected-text

  @Then /^no new files or folders have been created$/ (done) ->
    files-list = fs.readdir-sync path.join process.cwd!, 'tmp' |> unique
    new-files = difference files-list, @existing-folders
    if new-files.length > 0
      done new Error "New files or folders were created: #{new-files}"
    done!


