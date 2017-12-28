# Using dependency templates
- List dependency `<id>` as a remote dependency in `application.yml`, with one of the supported types
  - Populate required `template-config` fields (defined in the README for each type)
- List any dependency-specific service data under `dependency-data.<id>` in a service's `service.yml`

# Building dependency templates
- Create a directory with a unique name. This will be the dependency's `type`.
- The directory should contain the following files: `requirements.yml`, `dependency.tf`, a `README.md`, as well as a `modules` folder if applicable.
 Each is explained in detail below. (see existing templates for example)

## dependency.tf
A terraform file that contains the public interfacing terraform modules for a dependency. The file should be built
 as a [Mustache template](https://github.com/hoisie/mustache) which will be filled in with the user defined template-config.
 You can ensure the user passes in all the required fields using the requirements.yml file.
Example:
```tf
# dependency.tf
module "example_module" {
  version = "{{version}}"
}
```
```yml
# application.yml
remote:
  dependencies:
    <dependency-id>:
      type: example
      template-config:
        version: 0.0.1
```
```tf
# rendered dependency in user's main.tf
module "example_module" {
  version = "0.0.1"
}
```

### Terraform variables
The following are provided as Terraform variables
- `var.<dependency-id>_env_vars`: an environment variable passed to the dependency, containing a single key `SERVICE_DATA`, which contains a compilation
of all `dependency-data.<dependency-id>` listed in each `service.yml`
- `var.<secret-name>`: Any secrets created via `exo configure`
- `var.aws_profile`: Name of AWS profile passed into `exo deploy -p <profile-name>`
- `var.region`: AWS region set in `application.yml`
- `var.aws_account_id`: AWS account id set in `application.yml`
- `var.aws_ssl_certificate_arn`: ssl certificat arn set in `application.yml`
- `var.application_url`: applicaiton url set in `application.yml`
- `var.env`: deployment environment set in `exo deploy <env>`
- `module.aws.<output-variable>`: See `terraform/aws/vars.tf` for module output variables available for use

Submodules should be sourced using git URLs: https://www.terraform.io/docs/modules/sources.html#github (see below for more details).
The following fields are automatically rendered into the mustache template for each dependency:
- `terraformCommitHash`: the git commit hash of the latest changes to any terraform files of this repository

## modules
An optional directory containing submodules of `dependency.tf`.

## requirements.yml
Contains a single field: `required-fields`, which is a list of strings that define the fields that a user is required to populate
under the `template-config` section of a dependency configuration.

Example:
```yml
# requirements.yml

required-fields:
  - version
```
```yml
# application.yml

remote:
  dependencies:
    <dependency-id>:
      type: <type>
      template-config:
        version: 0.0.1
```

## README.md
Include a README with each dependency, be sure to include the following sections:
- Template configuration: describe each required field and provide sample values
- Environment variables: describe what environment variables the user is to set in their yaml files
- Service data: if there is any dependency-specific data for each service, describe what those objects should look like.
This would include any information listed under a `dependency-data.<dependency-id>` key in `application.yml` and `service.yml`
