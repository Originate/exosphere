require! {
  './deployers/aws-deployer' : AwsDeployer
}


# Deploys the overall application
class AppDeployer

  (@app-config) ->
    @deployer = new AwsDeployer @app-config


  deploy: (done) ->
    @deployer
      ..pull-remote-state ~>
        ..generate-terraform!
        ..deploy (err) ->
          | err => process.stdout.write "Error deploying application #{err.message}" ; return done err
          process.stdout.write "Application successfully deployed!"
          done!


  nuke: (done) ->
    @deployer
      ..pull-remote-state ~>
        ..nuke (err) ->
          | err => process.stdout.write "Error destroying application #{err.message}" ; return done err
          process.stdout.write "Application successfully destroyed!"
          done!


module.exports = AppDeployer
