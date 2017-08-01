require! {
  'cucumber': {After, Before, set-default-timeout, set-world-constructor}
  './world': World
}

set-default-timeout 1000
set-world-constructor World


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
