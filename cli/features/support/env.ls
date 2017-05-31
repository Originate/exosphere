require! {
  'child_process'
  '../../../exosphere-shared' : {kill-child-processes, DockerHelper}
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 2000


  @After tags: ['~@e2e'], (scenario, done) ->
    if @app-dir
      fs.remove-sync @app-dir
    kill-child-processes ->
      DockerHelper.remove-containers done


  #stop and remove all running docker containers
  @After tags: ['~@docker-cleanup'], timeout: 20_000, (scenario, done) ->
    DockerHelper.remove-containers done
