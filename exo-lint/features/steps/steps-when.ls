require! {
  'cucumber': {defineSupportCode}
  '../../../exosphere-shared' : {run-process}
  'path'
}


defineSupportCode ({When}) ->

  When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @app-dir
      ..on 'ended', -> done!
