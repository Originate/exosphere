require! {
  'chai' : {expect}
  'strip-ansi'
}


module.exports = ->

  @When /^running "([^"]*)" in the terminal$/, timeout: 10_000, (command, done) ->
    @run command
      ..on 'ended', (err, exit-code) ->
        done err


  @When /^trying to run "([^"]*)"$/ (command, done) ->
    @run command
      ..on 'ended', (err, exit-code) ~>
        expect(strip-ansi @process.full-output!).to.contain 'Error'
        done!
