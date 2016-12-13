require! {
  '../../../exosphere-shared' : {kill-child-processes}
  'rimraf'
}


module.exports = ->

  @set-default-timeout 2000


  @Before ->
    rimraf.sync 'tmp'


  @After (scenario, done) ->
    kill-child-processes done
