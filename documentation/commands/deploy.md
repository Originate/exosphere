# exo deploy

_Deploys an Exosphere application into the cloud_

Usage: `exo deploy [flags]`

Flags:
- `-p, --profile string`   AWS profile to use (default "default")

Deploys an application to the cloud, leveraging technology provided by [Terraform](https://terraform.io):
- Prepares AWS account for use with Terraform:
  - Creates S3 bucket to store Terraform state
  - Creates DynamoDB table to store Terraform lock
- Builds production Docker images and pushes them to Amazon's [EC2 Container Registry](https://aws.amazon.com/ecr/)
- Generates Terraform files based on application and service configuration
- Retrieves secrets managed by [`exo configure`](documentation/commands/configure.md) and passes them to Terraform processes
- Performs a dry run of deployment and outputs a plan of changes to be applied, asking for user confirmation
- Performs actual deployment

Currently supported platforms are AWS, more coming soon.
