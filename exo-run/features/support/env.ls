require! {
  'child_process'
  '../../../exosphere-shared' : {DockerHelper, kill-child-processes}
  'cucumber': {defineSupportCode}
  './world': World
}


defineSupportCode ({After, set-default-timeout, set-world-constructor}) ->

  set-default-timeout 2000

  set-world-constructor World

  After {timeout: 5000}, (scenario, done) ->
    kill-child-processes ->
      DockerHelper.remove-containers done
