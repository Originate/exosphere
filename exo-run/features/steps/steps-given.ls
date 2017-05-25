require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'observable-process' : ObservableProcess
  'path'
  'wait' : {wait}
}


module.exports = ->

  @Given /^a running "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      command = \exo-run
      if process.platform is 'win32' then command += '.cmd'
      @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                       cwd: @app-dir,
                                       stdout: dim-console.process.stdout
                                       stderr: dim-console.process.stderr)
      wait 5_000, done
