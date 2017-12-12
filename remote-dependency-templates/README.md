## Usage
- List dependency `<id>` as a remote dependency in `application.yml`, with one of the supported types
  - Populate required `template-config` fields (defined in the README for each type)
- List any dependency-specific service data under `dependency-data.<id>` in a service's `service.yml`

## Terraform variables //TODO

## Notes
The following are fields that exosphere passes to every dependency:
- `app-name`: Name of the application as defined in `application.yml`

