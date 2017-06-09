require! {
  'chai' : {expect}
  '../../../exosphere-shared' : {run-process}
  'path'
}


module.exports = ->

  @When /^trying to run "([^"]*)"$/, timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @current-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.not.equal 0
        done!

  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @process = run-process path.join(process.cwd!, 'bin', command), @current-dir
      ..on 'ended', (exit-code) ->
        expect(exit-code).to.equal 0
        done!
