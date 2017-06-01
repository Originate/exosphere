require! {
  'path'
}


module.exports = ->

  @Given /^a running "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-dir, ~>
      command = \exo-run
      if process.platform is 'win32' then command += '.cmd'
      @run-app {command, online-text: 'all services online'}, done 


