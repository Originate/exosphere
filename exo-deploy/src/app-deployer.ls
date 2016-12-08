require! {
  './deployers/aws-deployer' : AwsDeployer
}


# Deploys the overall application
class AppDeployer

  (@app-config) ->
    @deployer = new AwsDeployer @app-config


  deploy: ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform!
        ..deploy!


  nuke: ->
    @deployer
      ..pull-remote-state ~>
        ..nuke!


module.exports = AppDeployer
