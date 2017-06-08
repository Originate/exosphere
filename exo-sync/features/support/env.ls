require! {
  'cucumber': {defineSupportCode}
  './world': World
  '../../../exosphere-shared' : {kill-child-processes}
  'rimraf'
}


defineSupportCode ({After, Before, set-default-timeout, set-world-constructor})

  set-default-timeout 2000
  set-world-constructor World

  Before ->
    rimraf.sync 'tmp'


  After (scenario, done) ->
    kill-child-processes done
