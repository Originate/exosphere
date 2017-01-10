require! {
  'chai' : {expect}
  'child_process'
  'dim-console'
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name


  @Given /^I am in the root directory of an empty application called "([^"]*)" with the file "([^"]*)":$/, timeout: 10_000, (app-name, filename, file-content, done) !->
    @app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!
    fs.write-file-sync path.join(@app-dir, filename), file-content


  @Given /^The origin of "([^"]*)" contains a new commit not yet present in the local clone$/, (service-role) ->
    service-dir = path.join(@app-dir, service-type)
    service-origin-dir = @create-origin service-type, service-dir

    # create a new commit in the origin
    fs.write-file-sync path.join(service-origin-dir, 'new_file'), 'content'
    child_process.exec-sync "git add --all", cwd: service-origin-dir
    child_process.exec-sync "git status", cwd: service-origin-dir
    child_process.exec-sync "git commit -m new_commit", cwd: service-origin-dir



  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: path.join(@app-dir),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



  @When /^running "([^"]*)" in this application's "app" directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: path.join(@app-dir, 'app'),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^my application contains the newly committed file$/, ->
    fs.stat-sync path.join(@app-dir, 'web-service', 'new_file')
