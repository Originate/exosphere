require! {
  'cucumber': {After, Before, set-default-timeout}
}

set-default-timeout 1000


After (scenarioResult, done) ->
  closeIfDefined @exo-relay, ~>
    closeIfDefined @exocom, done


Before ->
  @ran = no


closeIfDefined = (obj, done) ->
  if obj
    obj.close done
  else
    done!
