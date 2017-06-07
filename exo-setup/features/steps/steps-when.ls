require! {
  'chai' : {expect}
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.equal 0
        done!

  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: process.stdout
                                     stderr: process.stderr)
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.equal 0
        done!
