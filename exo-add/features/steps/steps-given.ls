require! {
  'fs-extra' : fs
  'path'
  'yaml-cutter'
}


module.exports = ->

  # Note: The timeout exists because emptying the tmp dir might take a while.
  #       This is because the node_modules folder in there can contain a lot of files.
  @Given /^I am in the root directory of an empty application called "([^"]*)"$/, timeout: 20_000, (app-name, done) !->
    @app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app app-name, done


  @Given /^I am in an empty folder$/, ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir


  @Given /^I am in the directory of an application containing a "([^"]*)" service$/, (service-role, done) !->
    @app-dir := path.join process.cwd!, 'tmp', 'app'
    @create-empty-app 'app', ~>
      options =
        file: path.join @app-dir, 'application.yml'
        root: 'services.public'
        key: service-role
        value: {location: "./#{service-role}"}
      yaml-cutter.insert-hash options, done
