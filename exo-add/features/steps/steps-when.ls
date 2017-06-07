require! {
  'async'
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  @When /^running "([^"]*)" in the terminal$/, (command, done) ->
    @app-dir := path.join process.cwd!, 'tmp'
    @run command
      ..on 'ended', done


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    @run command
      ..on 'ended', done


  @When /^trying to run "([^"]*)" in this application's directory$/, timeout: 10_000, (command, done) ->
    @run command
      ..on 'ended', ~>
        | @process.exit-code > 0  =>  done!
        | otherwise               =>  throw new Error "Expected failure but exited with code 0"


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    @run command


  @When /^waiting until the process ends$/, timeout: 300_000, (done) ->
    @process.on 'ended', done
