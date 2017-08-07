require! {
  'exosphere-shared' : {kill-child-processes}
}


module.exports = ->

  @set-default-timeout 1500


  @After (scenario, done) ->
    @server1?.close!
    @server2?.close!
    @process?.kill!
    @exocom?.close!
    @exoservice?.close!
    kill-child-processes done
