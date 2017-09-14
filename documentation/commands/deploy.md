# exo deploy

_Deploys an Exosphere application into the cloud_

Usage: `exo deploy [flags]`

Flags:
- `-p, --profile string`   AWS profile to use (defaults to "default")

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
- Install both the AWS and Terraform command line interfaces

  Using Homebrew:
  - `brew install awscli`
  - `brew install terraform` (`>= 0.10.0` required)

  Or:
  - [AWS CLI installation](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
  - [Terraform CLI installation](https://www.terraform.io/intro/getting-started/install.html)
- Setup an [AWS account](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
- Configure [AWS CLI credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) on your machine
- Create and approve an [SSL certificate](https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-request.html) for your domain


#### Production configuration
- In `application.yml`, the following production fields are required:
  - `url`: apex domain application will use
  - `account-id`: AWS account number
  - `region`: AWS region to deploy resources to
  - `ssl-certificate-arn`: AWS [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of ssl certificate (see above)

Example:
```
production:
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
  - `public-port`: Service container port to expose
  - `health-check`: Endpoint where AWS will hit to perform health checks

  For a private service, the following fields are required:
  - `cpu`: see above
  - `memory`: see above
  - `public-port`: see above
  - `health-check`: see above

  For a worker service, the following fields are required:
  - `cpu`: see above
  - `memory`: see above

Example for a public service:
```
production:
  url: example.com
  cpu: 128
  memory: 128
  public-port: 3000
  health-check: '/'
```

#### Service environment variables
- Add public production environment variables to `environment/production` in each service's `service.yml`:
```
environment:
  default:
  development:
  production:
    ENV_VAR_NAME1: ENV_VAR_VALUE1
    ENV_VAR_NAME2: ENV_VAR_VALUE2
  secrets:
```
- Add private production environment variables to `environment/secrets` in each service's `service.yml` (see [exo configure](documentation/commands/configure.md))

#### Dependencies
Any dependencies that are to be ignored in production, or for which a third-party solution is desired,
 set the `external-in-production` field to be `true` under `dependencies/{dependency_name}/config/external-in-production` in either `service.yml` or `application.yml`.
```
dependencies:
  - name: mongo
    version: 3.4.0
    config:
      external-in-production: true
```

#### Optional Terraform variables
- `key_name` is the name of an EC2 Key Pair used to SSH into cloud instances.
Follow instructions for [Creating a Key Pair](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html?icmpid=docs_ec2_console) and create a secret `key_name = #{key_pair_name}` using `exo configure`. This key pair name will deployed with the machines.

### Service types
- Public: A service with an external facing [Application Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) that can accept external traffic
- Private: A service with an interal facing ALB, closed to external traffic
- Worker: A service with no ALBs, closed to external traffic

Currently supported platforms are AWS, more coming soon.
