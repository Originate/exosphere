require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!
