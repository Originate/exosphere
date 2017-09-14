# exo configure

_Manages an Exosphere application's secrets_

Usage:
- `exo configure [flags]` Creates a S3 bucket to store application secrets

Available subcommands:
- `exo configure create` creates secret key entries in remote secrets store
- `exo configure delete` deletes secrets from the remote secrets store
- `exo configure read` reads and prints secrets from remote secrets store
- `exo configure update` updates secret key entries in remote secrets store

Flags:
- `-p, --profile string`   AWS profile to use (default "default")

`exo configure` can be used to share application sensitive information between developers.
 It is also how `exo deploy` interpolates secret variables into Terraform files during deployment (see (documentation/commands/deploy.md)).

#### use with `exo deploy`
 To create a secret to be used in deployment, first define its name under `environment/secrets` of a service's `service.yml` file:

```
environment:
  default:
  development:
  production:
  secrets:
    - MONGODB_USER
    - MONGODB_PW
```

Second, create/manage secrets using the subcommands described above. Secrets are stored on S3 in the format: `secret_key = secret_value`.
 `secret_key` must match the corresponding secret name listed under `environment/secrets` in `service.yml`. The value of `secret_value` is
 the string that Terraform injects into the service during deployment.
