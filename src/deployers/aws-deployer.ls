require! {
  '../terraform/terraform' : Terraform
  '../terraform/aws-terraform-file-builder' : AwsTerraformFileBuilder
}


# Deploys application to AWS
class AwsDeployer

  (@app-config, @logger) ->
    process.env.AWS_ACCESS_KEY ? throw new Error "AWS access key not provided"
    process.env.AWS_SECRET_KEY ? throw new Error "AWS secret key not provided"
    @aws-config = @app-config.environments.production.providers.aws
    @exocom-port = 3100
    @exocom-dns = "exocom.#{@aws-config.region}.aws.#{@app-config.environments.production.domain}"
    @terraform = new Terraform


  generate-terraform: (done) ->
    new AwsTerraformFileBuilder {@app-config, @exocom-port, @exocom-dns}
      ..generate-terraform done @logger.log name: 'exo-deploy', text: "terraform generated for AWS"


  pull-remote-state: (done) ->
    bucket-path = @aws-config['remote-state-store'] |> (.split '/')
    unless bucket-path.length > 1 then throw new Error "application.yml param 'remote-state-store' missing S3 bucket and folder"
    backend-config = [
      "bucket=#{bucket-path[0]}"
      "key=#{bucket-path[1]}/terraform.tfstate"
      "region=#{@aws-config.region}"
      "access_key=#{process.env.AWS_ACCESS_KEY}"
      "secret_key=#{process.env.AWS_SECRET_KEY}"
    ]

    @terraform.pull-remote-state {backend: 's3', backend-config} done @logger.log name: 'exo-deploy', text: "terraform remote state pulled"


  deploy: ->
    @terraform
      ..get ~>
        @logger.log name: 'exo-deploy', text: "terraform starting deploy to AWS"
        ..apply!


module.exports = AwsDeployer
