require! {
  'cucumber': {defineSupportCode}
  './world': World
  '../../../exosphere-shared' : {kill-child-processes}
}


defineSupportCode ({After, set-default-timeout, set-world-constructor}) ->

  set-default-timeout 2000
  set-world-constructor World


  After (scenario, done) ->
    kill-child-processes done

