require! {
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 2000


  @After tags: ['~@e2e'], ->
    if @app-dir
      fs.remove-sync @app-dir
