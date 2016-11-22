require! {
  './deployers/aws-deployer' : AwsDeployer
  'events' : {EventEmitter}
  'exosphere-shared' : {Logger}
}


# Deploys the overall application
class AppDeployer extends EventEmitter

  (@app-config, @logger) ->
    @deployer = new AwsDeployer @app-config


  start: ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform! #TODO: logger print
        ..deploy!


module.exports = AppDeployer
