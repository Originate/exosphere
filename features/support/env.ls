require! {
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 2000


  @After ->
    if @app-dir
      fs.remove-sync @app-dir
