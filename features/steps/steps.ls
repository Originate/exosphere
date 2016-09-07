require! {
  'chai' : {expect}
  'child_process'
  'dim-console'
  'fs'
  'observable-process' : ObservableProcess
  'path'
  'yaml-cutter'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name


  @Given /^I am in the root directory of an application called "([^"]*)" using an external service "([^"]*)"$/, timeout: 10_000, (app-name, service-name, done) !->
    app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then ->
      # Insert the new service into the templated application.yml
      options =
        file: path.join app-dir, 'application.yml'
        root: 'services'
        key: service-name
        value: {location: "../#{app-name}/#{service-name}"}
      yaml-cutter.insert-hash options, done


  @Given /^I am in the root directory of an empty application called "([^"]*)" with the file "([^"]*)":$/, timeout: 10_000, (app-name, filename, file-content, done) !->
    app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!
    fs.write-file-sync path.join(app-dir, filename), file-content


  @Given /^The origin of "([^"]*)" contains a new commit not yet present in the local clone$/, (repo-name, done) ->
    @create-repo repo-name
    repo-dir = path.join(process.cwd!, 'tmp' ,'repos', repo-name)
    child_process.exec-sync "git clone ../repos/#{repo-name}", cwd: app-dir
    child_process.exec-sync "git add --all", cwd: repo-dir
    child_process.exec-sync "git commit -m message", cwd: repo-dir
    done!



  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!



  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^my application contains the newly committed file "([^"]*)"$/, (file-path) ->
    fs.stat-sync path.join(app-dir, file-path)
