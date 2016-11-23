require! {
  'async'
  'chai' : {expect}
  'dim-console'
  'exosphere-shared' : {call-args}
  'fs-extra' : fs
  'jsdiff-console'
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  # Note: The timeout exists because emptying the tmp dir might take a while.
  #       This is because the node_modules folder in there can contain a lot of files.
  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, timeout: 20_000, (app-name, done) !->
    app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!


  @Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir


  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  @When /^running "([^"]*)" in the terminal$/, (command, done) ->
    app-dir := path.join process.cwd!, 'tmp'
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^waiting until the process ends$/, timeout: 300_000, (done) ->
    @process.on 'ended', done


  @Then /^my application contains the file "([^"]*)" with the content:$/, (file-path, expected-content, done) ->
    fs.read-file path.join(app-dir, file-path), N (actual-content) ->
      jsdiff-console actual-content.to-string!trim!, expected-content.trim!, done


  @Then /^my application contains the file "([^"]*)"$/, (file-path) ->
    expect(fs.exists-sync path.join(app-dir, file-path)).to.be.true


  @Then /^my application contains the file "([^"]*)" containing the text:$/, (file-path, expected-fragment, done) ->
    fs.read-file path.join(app-dir, file-path), N (actual-content) ->
      expect(actual-content.to-string!).to.contain expected-fragment.trim!
      done!


  @Then /^it prints "([^"]*)" in the terminal$/, (expected-text) ->
    expect(@process.full-output!).to.contain expected-text

  @Then /^I see:$/ (expected-text) ->
     expect(@process.full-output!).to.contain expected-text

  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done
