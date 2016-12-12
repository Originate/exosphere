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
        ..deploy (err) ->
          | err => process.stdout.write "Error deploying application #{err.message}" ; return
          process.stdout.write "Application successfully deployed!"


  nuke: ->
    @deployer
      ..pull-remote-state ~>
        ..nuke (err) ->
          | err => process.stdout.write "Error destroying application #{err.message}" ; return
          process.stdout.write "Application successfully destroyed!"


module.exports = AppDeployer
