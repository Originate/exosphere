require! {
  '../../../exosphere-shared' : {DockerHelper, kill-child-processes}
}


module.exports = (done) ->

  @set-default-timeout 2000


  @After (scenario, done) ->
    kill-child-processes ->
      DockerHelper.remove-containers done
