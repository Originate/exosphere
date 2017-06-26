module.exports = ->

  @set-default-timeout 1000


  @After (_scenario, done) ->
    @service?.close!
    @exocom?.close done
