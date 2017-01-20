require! {
  'async'
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^executing "([^"]*)"$/, timeout: 20_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                 cwd: @app-dir,
                                 stdout: dim-console.process.stdout
                                 stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^starting "([^"]*)" in the terminal$/, (command) ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  @When /^running "([^"]*)" in the terminal$/, timeout: 20_000, (command, done) ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done
