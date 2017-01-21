require! {
  'child_process'
  '../../../exosphere-shared' : {kill-child-processes}
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 2000


  @After tags: ['~@e2e'], (scenario, done) ->
    if @app-dir
      fs.remove-sync @app-dir
    kill-child-processes done


  #stop and remove all running docker containers
  @After tags: ['~@docker-cleanup'], timeout: 20_000, (scenario, done) ->
    running-containers = child_process.exec-sync 'docker ps -q' |> (.to-string!)
    if running-containers
      child_process.exec 'docker rm $(docker stop $(docker ps -q))', done
    else
      done!
