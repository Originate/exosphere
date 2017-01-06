require! {
  'dim-console'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^running "([^"]*)"$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: @current-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!
