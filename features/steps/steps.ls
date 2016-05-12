require! {
  'async'
  'chai' : {expect}
  'dim-console'
  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'request'
  'tmp'
  'tmplconv'
}


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name


  @Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    @setup-app @app-name, done


  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, (app-name, done) !->
    @app-dir = path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join process.cwd!, 'templates', 'create-app'
    tmplconv.render(src-path, @app-dir, data: data).then ->
      done!



  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  @When /^installing it$/, timeout: 300_000, (done) ->
    @setup-app @app-name, done


  @When /^(trying to run|running) "([^"]*)" in the terminal$/, (attempt, command, done) ->
    @process = new ObservableProcess(path.join('bin', command),
                                     verbose: yes,
                                     console: dim-console)
    if attempt is 'trying to run'
      @process.wait 'Error', done
    else
      done!


  @When /^running "([^"]*)" in this application's directory$/, (command) ->
    @process = new ObservableProcess(path.join('..', 'bin', command),
                                     cwd: path.join(process.cwd!, 'tmp'),
                                     verbose: yes,
                                     console: dim-console)


  @When /^waiting until I see "([^"]*)"$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @When /^starting it$/, timeout: 10_000, (done) ->
    @start-app @app-name, done


  @When /^starting the "([^"]*)" application$/, (@app-name, done) ->
    @start-app @app-name, done



  @Then /^ExoCom uses this routing:$/, (table, done) ->
    expected-routes = {}
    for row in table.hashes!
      expected-routes[row.COMMAND] or= {}
      for receiver in row.RECEIVERS.split(', ')
        (expected-routes[row.COMMAND].receivers or= []).push host: 'localhost', name: receiver
    request "http://localhost:8000/config.json", N (response, body) ->
      expect(response.status-code).to.equal 200
      actual-routes = JSON.parse(body).routes
      for _, data of actual-routes
        for receiver in data.receivers
          expect(receiver.port).to.be.at.least 3000
          delete receiver.port
          delete receiver.internal-namespace
      jsdiff-console actual-routes, expected-routes, done


  @Then /^I see "([^"]*)"$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^it (?:creates|has created) the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@app-dir, row.SERVICE, row.FOLDER), fs.F_OK


  @Then /^my application contains the file "([^"]*)" containing the text:$/, (file-path, expected-fragment, done) ->
    fs.readFile path.join('tmp', file-path), N (actual-content) ->
      expect(actual-content.to-string!).to.contain expected-fragment.trim!
      done!


  @Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.readFile path.join('tmp', file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  @Then /^my machine is running ExoCom$/, (done) ->
    @process.wait 'exocom  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done


  @Then /^my workspace contains the file "([^"]*)" with content:$/, (path, expected-content, done) ->
    fs.readFile path, (err, actual-content) ->
      jsdiff-console actual-content.toString!trim!, expected-content.trim!, done
