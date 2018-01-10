# exo deploy

_Deploys an Exosphere application into the cloud_

Usage: `exo deploy [remote-environment-id] [flags]`

Flags:
- `-p, --profile string`  AWS profile to use (defaults to "default")
- `--auto-approve`        Deploys changes without prompting for approval

Deploys an application to the cloud, leveraging technology provided by [Terraform](https://terraform.io):
- Prepares AWS account for use with Terraform:
  - Creates S3 bucket to store Terraform state
  - Creates DynamoDB table to store Terraform lock
- Builds production Docker images and pushes them to Amazon's [EC2 Container Registry](https://aws.amazon.com/ecr/)
- Generates Terraform files based on application and service configuration
- Retrieves secrets managed by [`exo configure`](documentation/commands/configure.md) and passes them to Terraform processes
- Performs a dry run of deployment and outputs a plan of changes to be applied, asking for user confirmation
- Performs actual deployment

### User setup
A few steps are required of the user for a fully functional deployment:
- Setup an [AWS account](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
- Configure [AWS CLI credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html) on your machine
- Create and approve an [SSL certificate](https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-request.html) for your domain


#### Production configuration
- In `application.yml`, the following remote fields are required for each remote environment:
  - `url`: apex domain application will use
  - `account-id`: AWS account number
  - `region`: AWS region to deploy resources to
  - `ssl-certificate-arn`: AWS [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of ssl certificate (see above)

Example:
```
remote:
  environments:
    <remote-environment-id>:
      url: example.com
      account-id: 12345678
      region: us-west-2
      ssl-certificate-arn: certificate_arn
```

- The `service.yml` production fields vary dependeing on service type (see below)

  For a public service, the following fields are required:
  - `url`: URL to hit service at
  - `cpu`: Number of CPU units to reserve for service container (see "cpu" under [Container Definitions/Environment](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#container_definition_environment))
  - `memory`: The hard limit (in MiB) of memory to allocate to the service container (see "memory" under [Container Definitions](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#container_definitions))
  - `health-check`: Endpoint where AWS will hit to perform health checks

  For a worker service, the following fields are required:
  - `cpu`
  - `memory`

Example for a public service:
```
remote:
  cpu: 128
  memory: 128
  environments:
    <remote-environment-id>:
      url: example.com
      health-check: '/'
```

#### Service environment variables
- Add public production environment variables to `remote.environments.#{remote-environment-id}.environment-variables` in each service's `service.yml`:
- Add private production environment variables to `remote.environments.#{remote-environment-id}.secrets` in each service's `service.yml` (see [exo configure](configure.md))
```
remote:
  environments:
    <remote-environment-id>:
      environment-variables:
        ENV_VAR_NAME1: ENV_VAR_VALUE1
        ENV_VAR_NAME2: ENV_VAR_VALUE2
      secrets:
        - SECRET1
```

#### Dependencies
Define any dependency used in deployment under the `remote.dependencies` field. See [remote dependency templates](remote-dependency-templates/README.md) for more details on required fields
```
remote:
  dependencies:
    exocom:
      type: exocom
      template-config:
        version: 0.27.0
```

#### Optional Terraform variables
- `key_name` is the name of an EC2 Key Pair used to SSH into cloud instances.
Follow instructions for [Creating a Key Pair](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html?icmpid=docs_ec2_console) and create a secret `key_name = #{key_pair_name}` using `exo configure`. This key pair name will deployed with the machines.

#### Debugging

Once your public keys are on the bastion instances, you can ssh into the EC2 boxes running the services via: `ssh -o ProxyCommand='ssh -W %h:%p ubuntu@#{bastion_public_ip}' ec2-user@#{ec2_private_ip}`

Currently supported platforms are AWS, more coming soon.
