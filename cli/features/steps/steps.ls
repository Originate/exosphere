require! {
  'async'
  'chai' : {expect}
  'cucumber': {defineSupportCode}
  'fs-extra' : fs
  '../../../exosphere-shared' : {DockerHelper}
  'jsdiff-console'
  'nitroglycerin' : N
  \observable-process : ObservableProcess
  'path'
  'ps-tree'
  'zombie' : Browser
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


defineSupportCode ({Given, When, Then}) ->

  Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-name, done


  Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir


  Given /^I cd into "([^"]*)"$/ (dir-name) ->
    @app-name = dir-name
    app-dir := path.join process.cwd!, 'tmp', dir-name


  Given /^the file "([^"]*)":$/ (filename, file-content) ->
    # Note: uncomment this for running later scenarios of "features/tutorial.feature"
    #       by themselves.
    # app-dir := path.join process.cwd!, 'tmp', 'todo-app'
    fs.write-file-sync path.join(app-dir, filename), file-content


  Given /^my application contains the template folder "([^"]*)"$/ (template-dir) ->
    template-name = template-dir.split("/")[1]
    @checkout-service-template @app-name, template-name


  When /^adding a todo entry called "([^"]*)" via the web application$/ (entry, done) ->
    @browser.visit 'http://localhost:3000/', N ~>
      @browser.fill 'input[name=text]', entry
              .press-button 'add todo', done


  When /^entering into the wizard:$/, timeout: 10_000, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  # Note: This sometimes runs with the "tmp" directory populated with a ton of files.
  #       Cleaning them up can some time.
  #       Hence the larger timeout here.
  When /^running "([^"]*)" in the terminal$/, timeout: 30_000, (command, done) ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @run command, app-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.be.falsy
        done!


  When /^starting "([^"]*)" in the terminal$/, timeout: 20_000, (command) ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir
    @run command, app-dir


  When /^starting "([^"]*)" in this application's directory$/, (command) ->
    @run command, app-dir


  When /^executing the abbreviated command "([^"]*)" in the terminal$/, (command) ->
    @run command, app-dir


  When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @run command, app-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.be.falsy
        done!


  When /^(?:waiting until )?I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done


  When /^waiting until the process ends$/, timeout: 300_000, (done) ->
    @process.on 'ended', done


  Then /^I stop all running processes$/, (done) ->
    ps-tree @process.process.pid, N (children) ~>
      for child in children
        process.kill child.PID
      process.kill @process.process.pid
      done!


  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it prints the following in the terminal:$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it does not print "([^"]*)" in the terminal$/, (unexpected-text) ->
    expect(@process.full-output!).to.not.contain unexpected-text


  Then /^the docker images have the following folders:$/, (table) ->
    for row in table.hashes!
      process = new ObservableProcess(["docker", "run", row.IMAGE, "ls"]
                          cwd: @app-dir)
      process.on 'ended', ~>
        expect(process.full-output!).to.contain row.FOLDER


  Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.read-file path.join(app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  Then /^the full command "([^"]*)" is executed$/ timeout: 60_000, (command, done) ->
    expected-text = switch command
      | 'exo run'                => 'exo-run'
      | 'exo test'               => 'exo-test'
      | 'exo setup'              => 'exo-setup'
      | 'exo clean'              => 'We are about to clean up your Docker workspace'
      | 'exo create'             => 'We are about to create a new Exosphere application'
      | 'exo add'                => 'We are about to add a new Exosphere service to the application'
      | 'exo template'           => 'Manages remote service templates'
    @process.wait expected-text, done


  Then /^my workspace contains the file "([^"]*)" with content:$/, (filename, expected-content, done) ->
    fs.read-file path.join(app-dir, filename), N (actual-content) ->
      jsdiff-console actual-content.toString!trim!, expected-content.trim!, done


  Then /^my workspace contains the empty directory "([^"]*)"$/, (directory, done) ->
    fs.stat path.join(app-dir, directory), (err) ~>
      expect(err).to.be.null
      done!


  Then /^http:\/\/localhost:3000 displays:$/, timeout: 5_000, (expected-content, done) ->
    @browser or= new Browser
    @browser.visit 'http://localhost:3000/', N ~>
      @browser.assert.success!
      expect(@browser.text 'body').to.include expected-content.replace(/\n/g, '')
      done!
