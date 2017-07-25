# exo deploy

_Deploys an Exosphere application into the cloud_

Usage: `exo deploy`

Available subcommands:
- `exo deploy generate` generates a Terraform file
- `exo deploy run` runs the Terraform file

Deployment infrastructure is defined in the [Terraform](https://terraform.io) file
of the application.
Currently supported platforms are AWS, more coming soon.
