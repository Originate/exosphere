require! {
  'filendir' : {write-file}
  'fs-extra' : fs
  'mkdirp'
  'path'
}


module.exports = ->

  @Given /^I am in an empty folder$/, ->
    @app-dir := path.join process.cwd!, 'tmp'
    fs.empty-dir-sync @app-dir
    @existing-folders = []


  @Given /^source control contains a "([^"]*)" service$/ (service-name, done) ->
    repo-dir = path.join process.cwd!, 'tmp', 'origins', service-name
    mkdirp repo-dir, (err) ~>
      | err  =>  return done err
      fs.write-file-sync path.join(repo-dir, 'service.yml'), ' '
      @make-repo repo-dir
      done!


  @Given /^source control contains a repo "([^"]*)" with a file "([^"]*)" and the content:$/ (app-name, file-name, file-content, done) ->
    @existing-folders.push('origins')
    repo-dir = path.join process.cwd!, 'tmp', 'origins', app-name
    write-file path.join(repo-dir, file-name), file-content, (err) ~>
      @make-repo repo-dir
      done err


  @Given /^source control contains the services "([^"]*)" and "([^"]*)"$/ (service1, service2, done) ->
    @existing-folders.push('origins')
    repo-dirs = [ path.join(process.cwd!, 'tmp', 'origins', service1),
                  path.join(process.cwd!, 'tmp', 'origins', service2) ]
    for dir in repo-dirs
      mkdirp.sync dir
      file-name = path.join dir, 'service.yml'
      fs.write-file-sync file-name, ''
      @make-repo dir
    done!


  @Given /^my workspace already contains the folder "([^"]*)"$/ (folder-name, done) ->
    @existing-folders.push(folder-name)
    temp-file = path.join 'tmp', folder-name, 'temp-file'
    write-file temp-file, ' ', done
