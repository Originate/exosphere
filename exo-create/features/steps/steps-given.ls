require! {
  'path'
}


module.exports = ->

  # Note: The timeout exists because emptying the tmp dir might take a while.
  #       This is because the node_modules folder in there can contain a lot of files.
  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, timeout: 20_000, (app-name, done) !->
    @app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!

