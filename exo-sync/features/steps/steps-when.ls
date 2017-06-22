require! {
  'cucumber': {defineSupportCode}
  '../../../exosphere-shared' : {run-process}
  'path'
}


defineSupportCode ({When}) ->

  When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @app-dir
      ..on 'ended', done



  When /^running "([^"]*)" in this application's "app" directory$/, timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), path.join(@app-dir, 'app')
      ..on 'ended', done
