require! {
  'async'
  'chai' : {expect}
  'dim-console'
  'exosphere-shared' : {call-args}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'yaml-cutter'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^I am in the directory of an application with the services:$/, (table) ->
    app-dir := path.join process.cwd!, 'tmp', 'app'
    @create-empty-app('app')
      .then ->
        tasks = for row in table.hashes!
          options =
            file: path.join app-dir, 'application.yml'
            root: 'services'
            key: row.NAME
            value: {location: "./#{row.NAME}"}
          yaml-cutter.insert-hash(options, _)
        async.series tasks
      .then ~>
        @write-services table, app-dir


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text
