require! {
  'aws-sdk' : Aws
  'prelude-ls' : { map }
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
    @exocom-dns = "exocom.#{@aws-config.region}.aws.#{@domain-name}" #TODO: remove 'aws'
    @terraform = new Terraform
    @terraform-file-builder = new AwsTerraformFileBuilder {@app-config, @exocom-port, @exocom-dns}


  generate-terraform: (done) ->
    @_get-hosted-zone-id (hosted-zone-id) ~>
      @terraform-file-builder
        ..hosted-zone-id = hosted-zone-id
        ..generate-terraform done process.stdout.write "terraform scripts generated for AWS"


  pull-remote-state: (done) ->
    backend-config = [
      "bucket=#{@aws-config['remote-state-store']}"
      "key=terraform.tfstate"
      "region=#{@aws-config.region}"
    ]

    @_verify-remote-store ~>
      @terraform.pull-remote-state {backend: 's3', backend-config}, (err) ->
        | err => return process.stdout.write err.message
        process.stdout.write "terraform remote state pulled"
        done!


  deploy: ->
    @terraform
      ..get (err) ->
        | err => return process.stdout.write err.message
        process.stdout.write "terraform starting deploy to AWS"
        ..apply (err) ->
          | err => return process.stdout.write err.message


  teardown: ({nuke}) ->
    @terraform-file-builder.generate-provider-credentials!
    process.stdout.write "terraform starting nuke from AWS"
    @terraform.destroy (err) ~>
        | err => return process.stdout.write err.message
        if nuke then @_remove-hosted-zone @domain-name


  _get-hosted-zone-id: (done) ->
    @_hosted-zone-exists @domain-name, (id) ~>
      if id then done id
      else @_create-hosted-zone @domain-name, (id) -> done id


  _hosted-zone-exists: (domain-name, done) ->
    @route53 = new Aws.Route53 {api-version: '2013-04-01'}
      ..list-hosted-zones null, (err, data) ~>
        | err => return process.stdout.write err.message
        for hosted-zone in data.HostedZones
          if hosted-zone.Name is "#{domain-name}." then return done hosted-zone.Id
        return done no


  _create-hosted-zone: (domain-name, done) ->
    params =
      CallerReference: uuid.v4!
      Name: domain-name
    @route53.create-hosted-zone params, (err, data) ~>
      | err => return process.stdout.write err.message
      process.stdout.write "Please add the following name servers to #{@domain-name}:\n"
      for name-server in data.DelegationSet.NameServers
        process.stdout.write "#{name-server}\n"
      done data.HostedZone.Id


  _remove-hosted-zone: (domain-name) ->
    @_hosted-zone-exists @domain-name, (id) ~>
      if id then @route53.delete-hosted-zone {Id: id}, (err) ->
                    | err => return process.stdout.write err.message


  _verify-remote-store: (done) ~>
    @s3 = new Aws.S3 do
      api-version: '2006-03-01'
      region: @aws-config.region
    @_has-bucket @aws-config['remote-state-store'], (has-bucket) ~>
      unless has-bucket
        then @_create-remote-store done
        else done!


  _create-remote-store: (done) ->
    @s3
      ..create-bucket Bucket: @aws-config['remote-state-store'], CreateBucketConfiguration: LocationConstraint: "#{@aws-config.region}", (err, data) ~>
          if err then return done new Error err #TODO: inject logger object
          ..put-bucket-versioning Bucket: @aws-config['remote-state-store'], VersioningConfiguration: {Status: 'Enabled'}, ~>
            ..put-object Bucket: @aws-config['remote-state-store'], Key: 'terraform.tfstate', done


  # verify that s3 bucket with bucket-name exists
  _has-bucket: (bucket-name, done) ->
    @s3.list-buckets (err, data) ~>
      done bucket-name in @_bucket-names data


  # parses bucket names from Aws.s3.list-buckets callback
  _bucket-names: (data) ->
    data.Buckets |> map (.Name)


module.exports = AwsDeployer
