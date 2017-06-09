require! {
  'child_process'
  'cucumber': {defineSupportCode}
  'fs'
  'path'
}


defineSupportCode ({Given}) ->

  Given /^a freshly checked out "([^"]*)" application$/, (@app-name) ->
    @checkout-app @app-name
    @app-dir := path.join process.cwd!, 'tmp', @app-name


  Given /^I am in the root directory of an empty application called "([^"]*)" with the file "([^"]*)":$/, timeout: 10_000, (app-name, filename, file-content, done) !->
    @app-dir := path.join process.cwd!, 'tmp', app-name
    @create-empty-app(app-name, done).then -> done!
    fs.write-file-sync path.join(@app-dir, filename), file-content


  Given /^The origin of "([^"]*)" contains a new commit not yet present in the local clone$/, (service-name) ->
    service-dir = path.join(@app-dir, service-name)
    service-origin-dir = @create-origin service-name, service-dir

    # create a new commit in the origin
    fs.write-file-sync path.join(service-origin-dir, 'new_file'), 'content'
    child_process.exec-sync "git add --all", cwd: service-origin-dir
    child_process.exec-sync "git status", cwd: service-origin-dir
    child_process.exec-sync "git commit -m new_commit", cwd: service-origin-dir
