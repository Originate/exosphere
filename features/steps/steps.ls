require! {
  'async'
  'chai' : {expect}
  'child_process'
  'dim-console'
  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'ps-tree'
  'request'
  'tmp'
  'tmplconv'
  'yaml-cutter'
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


  @Given /^I am in the directory of the "([^"]*)" application$/ (@app-name) ->
    app-dir := path.join process.cwd!, 'example-apps', @app-name
    @current-dir = app-dir

  @Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir


  # Note: The timeout exists because emptying the tmp dir might take a while.
  #       This is because the node_modules folder in there can contain a lot of files.
  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, timeout: 20_000, (app-name, done) !->
    app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!


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


  @Given /^I cd into "([^"]*)"$/ (dir-name) ->
    app-dir := path.join process.cwd!, 'tmp', dir-name


  @Given /^the file "([^"]*)":$/ (filename, file-content) ->
    # Note: uncomment this for running later scenarios of "features/tutorial.feature"
    #       by themselves.
    # app-dir := path.join process.cwd!, 'tmp', 'todo-app'
    fs.write-file-sync path.join(app-dir, filename), file-content


  @Given /^The origin of "([^"]*)" contains a new commit not yet present in the local clone$/, (repo-name, done) ->
    @create-repo repo-name
    repo-dir = path.join(process.cwd!, 'tmp' ,'repos', repo-name)
    child_process.exec-sync "git clone ../repos/#{repo-name}", cwd: app-dir
    child_process.exec-sync "git add --all", cwd: repo-dir
    child_process.exec-sync "git commit -m message", cwd: repo-dir
    done!


  @Given /^source control contains the services "([^"]*)" and "([^"]*)"$/ (service1, service2) ->
    repo-dirs = [ path.join(process.cwd!, 'tmp', 'origins', service1),
                  path.join(process.cwd!, 'tmp', 'origins', service2) ]
    for dir in repo-dirs
      fs.mkdirs-sync dir
      file-name = path.join dir, 'service.yml'
      fd = fs.open-sync file-name, 'w'
      fs.write-file-sync file-name, ''
      fs.close-sync(fd)
      @make-repo dir


  @Given /^source control contains a repo "([^"]*)" with a file "([^"]*)" and the content:$/ (app-name, file-name, file-content) ->
    repo-dir = path.join process.cwd!, 'tmp', 'origins', app-name
    fs.mkdirs-sync repo-dir
    fs.write-file-sync path.join(repo-dir, file-name), file-content
    @make-repo repo-dir


  @Given /^I am in the "([^"]*)" directory$/ (service-dir, done) ->
    @checkout-app service-dir.split(path.sep)[0]
    @current-dir = path.join process.cwd!, 'tmp', service-dir
    done!


  @Given /^I am in the "([^"]*)" created directory$/ (dir, done) ->
    @checkout-app dir.split(path.sep)[0]
    new-dir = path.join process.cwd!, 'tmp', dir
    fs.mkdir-sync new-dir
    @current-dir = new-dir
    done!


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
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^running "([^"]*)" in the "([^"]*)" directory$/ (command, directory, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, directory)
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^running "([^"]*)" in directory ([^"]*)$/, timeout: 20_000, (command, dir, done) ->
    @process = new ObservableProcess(command,
                                     cwd: (path.join process.cwd!, dir),
                                     console: dim-console.console)
      ..on 'ended', done


  @When /^starting "([^"]*)" in the terminal$/, timeout: 20_000, (command) ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^executing the abbreviated command ([^"]*) in the terminal$/, (command) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                 cwd: app-dir,
                                 stdout: off
                                 stderr: off)


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^executing "([^"]*)"$/ (command, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                 cwd: app-dir,
                                 stdout: dim-console.process.stdout
                                 stderr: dim-console.process.stderr)
      ..on 'ended', done


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
    request "http://localhost:8001/config.json", N (response, body) ->
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


  @Then /^I stop all running processes$/, (done) ->
    ps-tree @process.process.pid, N (children) ~>
      for child in children
        process.kill child.PID
      process.kill @process.process.pid
      done!


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  @Then /^it (?:creates|has created) the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(app-dir, row.SERVICE, row.FOLDER), fs.F_OK


  @Then /^it creates the files:$/ (table) ->
    for row in table.hashes!
      fs.access-sync path.join(process.cwd!, 'tmp', row.FOLDER, row.FILE)


  @Then /^my application contains the file "([^"]*)"$/, (file-path) ->
    expect(fs.exists-sync path.join(app-dir, file-path)).to.be.true


  @Then /^my application contains the newly committed file "([^"]*)"$/, (file-path) ->
    fs.stat-sync path.join(app-dir, file-path)


  @Then /^my application contains the file "([^"]*)" containing the text:$/, (file-path, expected-fragment, done) ->
    fs.readFile path.join(app-dir, file-path), N (actual-content) ->
      expect(actual-content.to-string!).to.contain expected-fragment.trim!
      done!


  @Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.readFile path.join(app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  @Then /^my machine is running ExoCom$/, (done) ->
    @process.wait 'exocom  online at port', done

  @Then /^the full command "([^"]*)" is executed$/ (command, done) ->
    expected-text = switch command
      | 'exo run'                => 'exorun'
      | 'exo test'               => 'exo-test'
      | 'exo setup'              => 'exo-setup'
      | 'exo clone'              => 'We are going to clone an Exosphere application'
      | 'exo create application' => 'We are about to create a new Exosphere application'
      | 'exo create service'     => 'We are about to create a new Exosphere service'
      | 'exo add'                => 'We are about to add a new Exosphere service to the application'
    @process.wait expected-text, ~>
      @process
        ..kill!
        ..on 'ended', done


  @Then /^my machine is running the services:$/, (table, done) ->
    async.each [row['NAME'].to-lower-case! for row in table.hashes!],
               ((name, cb) ~> @process.wait "'#{name}' is running", cb),
               done


  @Then /^my workspace contains the file "([^"]*)" with content:$/, (filename, expected-content, done) ->
    fs.readFile path.join(app-dir, filename), N (actual-content) ->
      jsdiff-console actual-content.toString!trim!, expected-content.trim!, done


  @Then /^http:\/\/localhost:3000 displays:$/, timeout: 5_000, (expected-content, done) ->
    @browser or= new Browser
    @browser.visit 'http://localhost:3000/', N ~>
      @browser.assert.success!
      expect(@browser.text 'body').to.include expected-content.replace(/\n/g, '')
      done!


  @Then /^the "([^"]*)" service receives a "([^"]*)" message$/, (service, message, done) ->
    @process.wait "'#{service}' service received message '#{message}'", done


  @Then /^the "([^"]*)" service replies with a "([^"]*)" message$/, (arg1, arg2, done) ->
    done!


  @Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code


  @Then /^it only run tests for "([^"]*)"$/ (service-name, done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    @process.wait "exo-test  Testing service '#{service-name}'", done


  @Then /^it doesn't run any tests$/ (done) ->
    expect(@process.full-output!).to.not.include "Testing application"
    expect(@process.full-output!).to.not.include "Testing service"
    @process.wait "exo-test  Tests do not exist. Not in service or application directory.", done
