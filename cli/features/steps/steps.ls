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


  When /^executing the abbreviated command "([^"]*)" in the terminal$/, (command) ->
    @run command, app-dir


  Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it prints the following in the terminal:$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text


  Then /^it does not print "([^"]*)" in the terminal$/, (unexpected-text) ->
    expect(@process.full-output!).to.not.contain unexpected-text


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
