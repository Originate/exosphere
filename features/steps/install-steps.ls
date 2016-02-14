require! {
  '../support/dim-console'
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


module.exports = ->

  @When /^installing it$/, timeout: 60*1000, (done) ->
    @process = new ObservableProcess(path.join('..', '..', 'bin', 'exo-install'),
                                     cwd: path.join(process.cwd!, 'example-apps', @app-name),
                                     verbose: yes,
                                     console: dim-console)
      ..wait "installation complete", done


  @Then /^it creates the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@app-dir.name, row.SERVICE, row.FOLDER), fs.F_OK
