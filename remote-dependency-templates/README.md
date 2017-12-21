## Usage
- List dependency `<id>` as a remote dependency in `application.yml`, with one of the supported types
  - Populate required `template-config` fields (defined in the README for each type)
- List any dependency-specific service data under `dependency-data.<id>` in a service's `service.yml`

# Building dependency templates
- Create a directory with a unique name. This will be the dependency's `type`

## requirements.yml
`requirements.yml` contains a single field: `required-fields`, which is a list of strings that define the fields that a user is required to populate
under the `template-config` section of a dependency configuration.

Example:
```yml
# requirements.yml

required-fields:
  - version
```
```yml
# application.yml

<dependency-id>:
  type: <type>
  template-config:
    version: 0.0.1
```

## dependency.tf
Add a terraform file, `dependency.tf`, which contains the public interfacing terraform modules for a dependency. The `dependency.tf` file should be built
as a [Mustache template](https://github.com/hoisie/mustache). Any configuration information listed as a required field in `requirements.yml` (see above)
will be rendered into `dependency.tf` using Mustache. The mustache variables must match those listed in `requirements.yml`, and hence in a user's `template-config`.
Example:
```
# dependency.tf
module "example_module" {
  name = {{dependency-name}}
}
```
```yml
# requirements.yml
required-fields:
 - dependency-name
```
```
# application.yml
dependencies:
  <dependency-id>:
    type: example
    template-config:
      dependency-name: this-will-be-rendered
```

Submodules should be sourced using git URLs: https://www.terraform.io/docs/modules/sources.html#github (see below for more details).
The following fields are automatically rendered into the mustache template for each dependency:
- `terraformCommitHash`: The git commit hash of the most recent commit made to `terraform/aws`

## Terraform submodules
Any submodules of the main `dependency.tf` should be defined in `#{dependencyType}/modules`, and should be referenced from `dependency.tf` using git URLs.

## Terraform variables
The following are provided as Terraform variables, and can be used via `${var.<field-name>}` in `dependency.tf`:
- `env`: the deployment environment. i.e. qa, staging, production
- `key-name`: the EC2 key pair name to attach to any EC2 instances
- `<dependency-id>_env_vars`: an environment variable passed to the dependency, containing a single key `SERVICE_DATA`, which contains a compilation
of all `dependency-data.<dependency-id>` listed in each `service.yml`
- Any secrets created via `exo configure` will be available via `var.<secret-name>
- The application URL defined in `application.yml` will also be avialble via `var.application_url`, and each public service's url will also be available via `var.<service-role>_url`

## README
Include a README with each dependency, be sure to include the following sections:
- Template configuration: describe what `template-config` fields are required (those listed in `requirements.yml`)
- Environment variables: describe what environment variables the user is to set in their yaml files
- Service data: if there is any dependency-specific data for each service, describe what those objects should look like.
This would include any information listed under a `dependency-data.<dependency-id>` key in `application.yml` and `service.yml`
