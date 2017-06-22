require! {
  'cucumber': {defineSupportCode}
  'fs-extra' : fs
  './world': World
}


defineSupportCode ({After, set-default-timeout, set-world-constructor}) ->

  set-default-timeout 1000
  set-world-constructor World

  After ->
    @exocom?.close!
    @process?.close!
