require! {
  'child_process'
  '../../../exosphere-shared' : {kill-child-processes}
}


module.exports = (done) ->

  @set-default-timeout 2000


  @After (scenario, done) ->
    kill-child-processes ->
      if child_process.exec-sync 'docker ps -q'
        child_process.exec-sync 'docker stop $(docker ps -q)'
        child_process.exec-sync 'docker rm $(docker ps -aq)'
        done!
