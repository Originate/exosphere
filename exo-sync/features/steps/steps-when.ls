require! {
  'dim-console'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: path.join(@app-dir),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done



  @When /^running "([^"]*)" in this application's "app" directory$/, timeout: 600_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', command),
                                     cwd: path.join(@app-dir, 'app'),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done
