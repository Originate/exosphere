module.exports = ->

  @set-default-timeout 1000


  @After ->
    @exo-relay?.close!
    @exocom?.close!


  @Before ->
    @ran = no
