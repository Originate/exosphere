require! {
  'chai' : {expect}
  '../../../exosphere-shared' : {DockerHelper}
  'fs'
  'prelude-ls' : {map}
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
      # split image_name:version to get image_name only
      image-names = map (.split(':')[0]), docker-images
      console.log image-names
      for row in table.raw!
        expect(image-names).to.include row[0]
      done!


  @Then /^it finishes with exit code (\d+)$/ (+expected-exit-code) ->
    expect(@process.exit-code).to.equal expected-exit-code
