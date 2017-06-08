require! {
  'cucumber': {defineSupportCode} 
  'fs'
  'path'
}


defineSupportCode ({Given}) ->

  Given /^a set\-up "([^"]*)" application$/, timeout: 600_000, (@app-name, done) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name
    @setup-app @app-name, done


  Given /^I am in the "([^"]*)" directory$/ (service-dir, done) ->
    @checkout-app service-dir.split(path.sep)[0]
    @current-dir = path.join process.cwd!, 'tmp', service-dir
    done!


  Given /^I am in the "([^"]*)" created directory$/ (dir, done) ->
    new-dir = path.join process.cwd!, dir
    try
      fs.mkdir-sync new-dir
    @current-dir = new-dir
    done!
