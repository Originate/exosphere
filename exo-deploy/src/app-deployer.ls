require! {
  './deployers/aws-deployer' : AwsDeployer
}


# Deploys the overall application
class AppDeployer

  (@app-config, command-flag) ->
    @deployer = new AwsDeployer @app-config


  deploy: ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform!
        ..deploy!


  teardown: ({nuke}) ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform!
        if nuke then @deployer.nuke! else @deployer.teardown!


module.exports = AppDeployer
