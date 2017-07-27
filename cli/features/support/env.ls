require! {
  'cucumber': {defineSupportCode}
  './world': World
  'child_process'
  '../../../exosphere-shared' : {kill-child-processes, DockerHelper}
  'fs-extra' : fs
}


defineSupportCode ({After, set-default-timeout, set-world-constructor}) ->

  set-default-timeout 5000
  set-world-constructor World


  After tags: 'not @e2e', (scenario, done) ->
    if @app-dir
      fs.remove-sync @app-dir
    kill-child-processes ->
      DockerHelper.remove-containers done


  #stop and remove all running docker containers
  After tags: 'not @docker-cleanup', timeout: 20_000, (scenario, done) ->
    DockerHelper.remove-containers done
