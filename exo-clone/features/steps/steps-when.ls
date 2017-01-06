require! {
  'chai' : {expect}
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'strip-ansi'
}


module.exports = ->

  @When /^running "([^"]*)" in the terminal$/, timeout: 10_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, 'tmp')
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', (err, exit-code) ->
        done err


  @When /^trying to run "([^"]*)"$/ (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, 'tmp')
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', (err, exit-code) ~>
        expect(strip-ansi @process.full-output!).to.contain 'Error'
        done!
