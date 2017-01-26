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
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^running "([^"]*)" in this application's directory$/, timeout: 600_000, (command, done) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @When /^trying to run "([^"]*)" in this application's directory$/, timeout: 600_000, (command) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^starting "([^"]*)" in this application's directory$/, (command) ->
    args = command.split ' '
    args[0] = path.join process.cwd!, 'bin', args[0]
    if process.platform is 'win32'
      args[0] += '.cmd'
    @process = new ObservableProcess(args,
                                     cwd: @app-dir,
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)


  @When /^waiting until the process ends$/, timeout: 300_000, (done) ->
    @process.on 'ended', done

