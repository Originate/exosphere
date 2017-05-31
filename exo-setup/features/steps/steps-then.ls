require! {
  'chai' : {expect}
  '../../../exosphere-shared' : {DockerHelper}
  'fs'
  'jsdiff-console'
  'path'
}


module.exports = ->

  @Then /^it has created the folders:$/, (table) ->
    for row in table.hashes!
      fs.access-sync path.join(@current-dir, row.SERVICE, row.FOLDER)


  @Then /^it has created the files:$/, (table) ->
    for row in table.raw!
      fs.access-sync path.join(@current-dir, row[0])


  @Then /^it has acquired the Docker images:$/ (table, done) ->
    DockerHelper.list-images (err, docker-images) ->
      console.log docker-images
      for row in table.raw!
        expect(docker-images).to.include row[0]
      done!


  @Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code
