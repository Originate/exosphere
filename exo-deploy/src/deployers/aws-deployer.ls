require! {
  'aws-sdk' : Aws
  'prelude-ls' : { map }
  '../terraform/terraform' : Terraform
  '../terraform/aws-terraform-file-builder' : AwsTerraformFileBuilder
}


# Deploys application to AWS
class AwsDeployer

  (@app-config) ->
    process.env.AWS_ACCESS_KEY_ID ? throw new Error "AWS_ACCESS_KEY_ID not provided"
    process.env.AWS_SECRET_ACCESS_KEY ? throw new Error "AWS_SECRET_ACCESS_KEY not provided"
    @aws-config = @app-config.environments.production.providers.aws
    @exocom-port = 3100
    @exocom-dns = "exocom.#{@aws-config.region}.aws.#{@app-config.environments.production.domain}"
    @terraform = new Terraform


  generate-terraform: ->
    new AwsTerraformFileBuilder {@app-config, @exocom-port, @exocom-dns}
      ..generate-terraform process.stdout.write "terraform scripts generated for AWS"


  pull-remote-state: (done) ->
    backend-config = [
      "bucket=#{@aws-config['remote-state-store']}"
      "key=terraform.tfstate"
      "region=#{@aws-config.region}"
    ]

    @_verify-remote-store (err) ~>
      | err => return process.stdout.write "Cannot verify remote store #{err.message}"
      @terraform.pull-remote-state {backend: 's3', backend-config}, (err) ->
        | err => return process.stdout.write "Cannot pull remote state #{err.message}"
        process.stdout.write "terraform remote state pulled"
        done!


  deploy: ->
    @terraform
      ..get (err) ->
        | err => return process.stdout.write err.message
        process.stdout.write "terraform starting deploy to AWS"
        ..apply (err) ->
          | err => return process.stdout.write err.message


  _verify-remote-store: (done) ~>
    @s3 = new Aws.S3 do
      api-version: '2006-03-01'
      region: @aws-config.region
    @_has-bucket @aws-config['remote-state-store'], (err, has-bucket) ~>
      | err => process.stdout.write "Cannot verify remote store S3 bucket: #{err.message}" ; return done err
      if !has-bucket
        @_create-remote-store done, (err) ->
          | err => process.stdout.write "Cannot create remote store S3 bucket: #{err.message}" ; return done err
      else done!


  _create-remote-store: (done) ->
    @s3
      ..create-bucket Bucket: @aws-config['remote-state-store'], CreateBucketConfiguration: LocationConstraint: "#{@aws-config.region}", (err, data) ~>
          | err => process.stdout.write "Cannot create S3 bucket: #{err.message}" ; return done err
          ..put-bucket-versioning Bucket: @aws-config['remote-state-store'], VersioningConfiguration: {Status: 'Enabled'}, (err, data) ~>
            | err => process.stdout.write "Cannot configure remote store bucket versioning: #{err.message}" ; return done err
            ..put-object Bucket: @aws-config['remote-state-store'], Key: 'terraform.tfstate', (err, data) ->
                | err => process.stdout.write "Cannot put object terraform.tfstate: #{err.message}" ; return done err
                done null


  # verify that s3 bucket with bucket-name exists
  _has-bucket: (bucket-name, done) ->
    @s3.list-buckets (err, data) ~>
      | err => process.stdout.write "#{err.message}" ; return done err
      done null, bucket-name in @_bucket-names data


  # parses bucket names from Aws.s3.list-buckets callback
  _bucket-names: (data) ->
    data.Buckets |> map (.Name)


module.exports = AwsDeployer
