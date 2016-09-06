require! {
  'dim-console'
  'exosphere-shared' : {call-args}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
}


# We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null


module.exports = ->

  @Given /^I am in an empty folder$/, ->
    app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync app-dir


  @Given /^source control contains a repo "([^"]*)" with a file "([^"]*)" and the content:$/ (app-name, file-name, file-content) ->
    repo-dir = path.join process.cwd!, 'tmp', 'origins', app-name
    fs.mkdirs-sync repo-dir
    fs.write-file-sync path.join(repo-dir, file-name), file-content
    @make-repo repo-dir


  @Given /^source control contains the services "([^"]*)" and "([^"]*)"$/ (service1, service2) ->
    repo-dirs = [ path.join(process.cwd!, 'tmp', 'origins', service1),
                  path.join(process.cwd!, 'tmp', 'origins', service2) ]
    for dir in repo-dirs
      fs.mkdirs-sync dir
      file-name = path.join dir, 'service.yml'
      fd = fs.open-sync file-name, 'w'
      fs.write-file-sync file-name, ''
      fs.close-sync(fd)
      @make-repo dir


  @When /^running "([^"]*)" in the "([^"]*)" directory$/ (command, directory, done) ->
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     cwd: path.join(process.cwd!, directory)
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', (err, exit-code) ->
        console.log err
        console.log exit-code
        done!


  @Then /^it creates the files:$/ (table) ->
    for row in table.hashes!
      fs.access-sync path.join(process.cwd!, 'tmp', row.FOLDER, row.FILE)
