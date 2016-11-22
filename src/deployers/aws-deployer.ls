require! {
  '../terraform/terraform' : Terraform
  '../terraform/aws-terraform-file-builder' : AwsTerraformFileBuilder
}


# Deploys application to AWS
class AwsDeployer

  (@app-config) ->
    process.env.AWS_ACCESS_KEY ? throw new Error "AWS access key not provided"
    process.env.AWS_SECRET_KEY ? throw new Error "AWS secret key not provided"
    @exocom-port = 3100
    @exocom-dns = "exocom.#{@app-config.environments.production.providers.aws.region}.aws.#{@app-config.environments.production.domain}"


  generate-terraform: (done) ->
    new AwsTerraformFileBuilder {@app-config, @exocom-port, @exocom-dns}
      ..generate-terraform done


  deploy: ->
    new Terraform
      ..get ~>
        ..apply!


module.exports = AwsDeployer
