require! {
  '../../../exosphere-shared' : {kill-child-processes}
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 2000


  @After tags: ['~@e2e'], (scenario, done) ->
    if @app-dir
      fs.remove-sync @app-dir
    kill-child-processes done
