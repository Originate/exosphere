require! {
  'aws-sdk' : Aws
  'prelude-ls' : {find, map}
  '../terraform/terraform' : Terraform
  '../terraform/aws-terraform-file-builder' : AwsTerraformFileBuilder
  'uuid'
}


# Deploys application to AWS
class AwsDeployer

  (@app-config) ->
    process.env.AWS_ACCESS_KEY_ID ? throw new Error "AWS_ACCESS_KEY_ID not provided"
    process.env.AWS_SECRET_ACCESS_KEY ? throw new Error "AWS_SECRET_ACCESS_KEY not provided"
    @aws-config = @app-config.environments.production.providers.aws
    @exocom-port = 3100
    @domain-name = @app-config.environments.production.domain
    @exocom-dns = "exocom.#{@aws-config.region}.#{@domain-name}"
    @terraform = new Terraform
    @terraform-file-builder = new AwsTerraformFileBuilder {@app-config, @exocom-port, @exocom-dns}


  generate-terraform: ->
    @terraform-file-builder.generate-terraform process.stdout.write "terraform scripts generated for AWS"


  pull-remote-state: (done) ->
    backend-config = [
      "bucket=#{@aws-config['remote-state-store']}"
      "key=terraform.tfstate"
      "region=#{@aws-config.region}"
    ]

    @_verify-remote-store (err) ~>
      | err => return process.stdout.write "Cannot verify remote store #{err.message}"
      @terraform-file-builder.empty-terraform-dir!
      @terraform.pull-remote-state {backend: 's3', backend-config}, (err) ->
        | err => return process.stdout.write "Cannot pull remote state #{err.message}"
        process.stdout.write "terraform remote state pulled"
        done!


  deploy: (done) ->
    @terraform
      ..get (err) ~>
        | err =>  process.stdout.write err.message ; return done err
        @_get-hosted-zone-id (err, hosted-zone-id) ~>
          | err =>  process.stdout.write "Cannot get hosted zone id #{err.message}" ; return done err
          ..apply {hosted-zone-id}, (err) ->
            | err =>  process.stdout.write err.message ; return done err
            done!


  teardown: ({nuke}, done) ->
    @terraform-file-builder.generate-provider-credentials!
    process.stdout.write "terraform starting AWS teardown"
    @terraform.destroy (err) ~>
      | err =>  process.stdout.write "Terraform cannot destroy infrastructure #{err.message}" ; return done err
      if nuke
        @_remove-hosted-zone @domain-name, (err) ->
          | err =>  process.stdout.write "Cannot remove hosted zone #{err.message}" ; return done err
          done!
      else done!


  # gets hosted zone id if it exists, creates one if it doesn't
  _get-hosted-zone-id: (done) ->
    @_get_hosted_zone @domain-name, (err, hosted-zone) ~>
      | err         =>  process.stdout.write err.message ; return done err
      | hosted-zone =>  done null, hosted-zone.Id
      | _           =>  @_create-hosted-zone @domain-name, done


  _get_hosted_zone: (domain-name, done) ->
    @route53 = new Aws.Route53 api-version: '2013-04-01'
      ..list-hosted-zones null, (err, data) ~>
        | err  =>  process.stdout.write "Cannot list hosted zones #{err.message}"
        done err, (data.HostedZones or [] |> find (.Name is "#{domain-name}."))


  _create-hosted-zone: (domain-name, done) ->
    params =
      CallerReference: uuid.v4!
      Name: domain-name

    @route53.create-hosted-zone params, (err, data) ~>
      | err => process.stdout.write "Cannot create hosted zone #{err.message}" ; return done err
      process.stdout.write "Please add the following name servers to #{@domain-name}:\n"
      for name-server in data.DelegationSet.NameServers
        process.stdout.write "#{name-server}\n"
      done null, data.HostedZone.Id


  _remove-hosted-zone: (domain-name, done) ->
    @_get_hosted_zone @domain-name, (err, hosted-zone) ~>
      | err           =>  process.stdout.write err.message ; return done err
      | hosted-zone   =>  @route53.delete-hosted-zone {Id: hosted-zone.Id}, (err) ->
                            | err => process.stdout.write err.message
                            done err
      | _             =>  done! # no hosted zone exists to remove


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
