require! {
  'exosphere-shared' : {kill-child-processes}
}


module.exports = ->

  @set-default-timeout 1000


  @After (scenario, done) ->
    @existing-server?.close!
    @exocom?.close!
    @exocom?.clearPorts!
    @process?.kill!
    for _, mock of @service-mocks or {}
      mock.close!
    kill-child-processes done
