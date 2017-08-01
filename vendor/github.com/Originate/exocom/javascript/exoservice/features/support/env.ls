require! {
  'exosphere-shared' : {kill-child-processes}
}


module.exports = ->

  @set-default-timeout 1500


  @After (scenario, done) ->
    @server1?.close!
    @server2?.close!
    @process?.kill!
    closeIfDefined @exoservice, ~>
      closeIfDefined @exocom, ~>
        kill-child-processes done

closeIfDefined = (obj, done) ->
  if obj
    obj.close done
  else
    done!
