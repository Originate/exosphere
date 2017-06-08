require! {
  'async'
  '../../../exosphere-shared' : {call-args}
  'fs-extra' : fs
  'path'
}


module.exports = ->

  @When /^executing "([^"]*)"$/, timeout: 20_000, (command, done) ->
    @run command
      ..on 'ended', done


  @When /^starting "([^"]*)" in the terminal$/, (command) ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    @run command


  @When /^entering into the wizard:$/, (table, done) ->
    enter-input = ([text, input], cb) ~>
      <~ @process.wait text
      @process.stdin.write "#{input}\n"
      cb!
    async.each table.rows!, enter-input, done


  @When /^running "([^"]*)" in the terminal$/, timeout: 20_000, (command, done) ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    @run command
      ..on 'ended', done


  @When /^waiting until I see "([^"]*)" in the terminal$/, timeout: 300_000, (expected-text, done) ->
    @process.wait expected-text, done
