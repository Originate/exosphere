# Using dependency templates
- List dependency `<id>` as a remote dependency in `application.yml`, with one of the supported types
  - Populate required `template-config` fields (defined in the README for each type)
- List any dependency-specific service data under `dependency-data.<id>` in a service's `service.yml`

# Building dependency templates
- Create a directory with a unique name. This will be the dependency's `type`.
- The directory should contain the following three files: requirements.yml, dependency.tf and README.md. Each is explained in detail below. (see existing templates for example)

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

<dependency-id>:
  type: <type>
  template-config:
    version: 0.0.1
```

## dependency.tf
A terraform file that contains the public interfacing terraform modules for a dependency. The file should be built
as a [Mustache template](https://github.com/hoisie/mustache). Any configuration information listed as a required field in `requirements.yml` (see above)
will be rendered into `dependency.tf` using Mustache. The mustache variables must match those listed in `requirements.yml`, and hence in a user's `template-config`.
Example:
```tf
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
```yml
# application.yml
dependencies:
  <dependency-id>:
    type: example
    template-config:
      dependency-name: this-will-be-rendered
```

Submodules should be sourced using git URLs: https://www.terraform.io/docs/modules/sources.html#github (see below for more details).
The following fields are automatically rendered into the mustache template for each dependency:
- `terraformCommitHash`: the git commit hash of the latest changes to any terraform files of this repository

### Terraform submodules
Any submodules of the main `dependency.tf` should be defined in `#{dependencyType}/modules`, and should be referenced from `dependency.tf` using git URLs.

### Terraform variables
The following are provided as Terraform variables
- `var.<dependency-id>_env_vars`: an environment variable passed to the dependency, containing a single key `SERVICE_DATA`, which contains a compilation
of all `dependency-data.<dependency-id>` listed in each `service.yml`
- `var.<secret-name>`: Any secrets created via `exo configure`
- Any variables listed at the top of `src/terraform/templates/aws.tf`

## README
Include a README with each dependency, be sure to include the following sections:
- Template configuration: describe each required field and provide sample values
- Environment variables: describe what environment variables the user is to set in their yaml files
- Service data: if there is any dependency-specific data for each service, describe what those objects should look like.
This would include any information listed under a `dependency-data.<dependency-id>` key in `application.yml` and `service.yml`
