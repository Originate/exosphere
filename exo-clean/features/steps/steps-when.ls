require! {
  'cucumber': {defineSupportCode}
  '../../../exosphere-shared' : {run-process}
  'path'
}


defineSupportCode ({When}) ->

  When /^running "([^"]*)" in the terminal$/ timeout: 6_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command)
      ..on 'ended', -> done!
