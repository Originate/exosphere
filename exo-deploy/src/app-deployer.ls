require! {
  './deployers/aws-deployer' : AwsDeployer
}


# Deploys the overall application
class AppDeployer

  (@app-config) ->
    @deployer = new AwsDeployer @app-config


  start: ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform!
        ..deploy!


module.exports = AppDeployer
