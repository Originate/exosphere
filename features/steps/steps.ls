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
  'zombie' : Browser
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name


  @Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-name, done


  @Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir


  # Note: The timeout exists because emptying the tmp dir might take a while.
  #       This is because the node_modules folder in there can contain a lot of files.
  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, timeout: 10_000, (app-name, done) !->
    app-dir := path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join process.cwd!, 'templates', 'create-app'
    tmplconv.render(src-path, app-dir, {data}).then -> done!


  @Given /^I cd into "([^"]*)"$/ (dir-name) ->
    app-dir := path.join process.cwd!, 'tmp', dir-name


  @Given /^the file "([^"]*)":$/ (filename, file-content) ->
    # Note: uncomment this for running later scenarios of "features/tutorial.feature"
    #       by themselves.
    # app-dir := path.join process.cwd!, 'tmp', 'todo-app'
    fs.write-file-sync path.join(app-dir, filename), file-content


  @When /^adding a todo entry called "([^"]*)" via the web application$/ (entry, done) ->
    @browser.visit 'http://localhost:3000/', N ~>
      @browser.fill 'input[name=text]', entry
              .press-button 'add todo', done


  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  # Note: This sometimes runs with the "tmp" directory populated with a ton of files.
  #       Cleaning them up can some time.
  #       Hence the larger timeout here.
  @When /^running "([^"]*)" in the terminal$/, timeout: 20_000, (command, done) ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     console: dim-console.console)
      ..on 'ended', done


  @When /^starting "([^"]*)" in the terminal$/, timeout: 20_000, (command) ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     console: dim-console.console)


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     console: dim-console.console)

  @When /^running "([^"]*)" in this application's directory$/, timeout: 300_000, (command, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     console: dim-console.console)
      ..on 'ended', (exit-code) -> done!


  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  @When /^waiting until the process ends$/, timeout: 300_000, (done) ->
    @process.on 'ended', done


  @When /^starting it$/, timeout: 10_000, (done) ->
    @start-app @app-name, done


  @When /^starting the "([^"]*)" application$/, (@app-name, done) ->
    @start-app @app-name, done


  @When /^the web service broadcasts a "([^"]*)" message$/, (message, done) ->
    request 'http://localhost:4000', done



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


  @Then /^I kill the server$/, (done) ->
    @process
      ..on 'ended', -> done!
      ..kill!


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^it (?:creates|has created) the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(app-dir, row.SERVICE, row.FOLDER), fs.F_OK


  @Then /^my application contains the file "([^"]*)"$/, (file-path) ->
    expect(fs.exists-sync path.join(app-dir, file-path)).to.be.true



  @Then /^my application contains the file "([^"]*)" containing the text:$/, (file-path, expected-fragment, done) ->
    fs.readFile path.join(app-dir, file-path), N (actual-content) ->
      expect(actual-content.to-string!).to.contain expected-fragment.trim!
      done!


  @Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.readFile path.join(app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  @Then /^my machine is running ExoCom$/, (done) ->
    @process.wait 'exocom  online at port', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done


  @Then /^my workspace contains the file "([^"]*)" with content:$/, (filename, expected-content, done) ->
    fs.readFile path.join(app-dir, filename), N (actual-content) ->
      jsdiff-console actual-content.toString!trim!, expected-content.trim!, done


  @Then /^http:\/\/localhost:3000 displays:$/ (expected-content, done) ->
    @browser or= new Browser
    @browser.visit 'http://localhost:3000/', N ~>
      @browser.assert.success!
      expect(@browser.text 'body').to.include expected-content.replace(/\n/g, '')
      done!


  @Then /^the "([^"]*)" service receives a "([^"]*)" message$/, (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done


  @Then /^the "([^"]*)" service replies with a "([^"]*)" message$/, (arg1, arg2, done) ->
    done!
