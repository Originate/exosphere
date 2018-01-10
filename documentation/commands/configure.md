# exo configure

_Manages an Exosphere application's secrets_

Usage:
- `exo configure [command]` Creates a S3 bucket to store application secrets

Available commands:
- `exo configure create [remote-environment-id] [flags]` creates secret in remote secrets store
- `exo configure delete [remote-environment-id] [flags]` deletes secret from remote secrets store
- `exo configure read [remote-environment-id] [flags]` prints secrets from remote secrets store
- `exo configure update [remote-environment-id] [flags]` updates secret in remote secrets store

Flags:
- `-p, --profile string`   AWS profile to use (defaults to "default")

`exo configure` can be used to share application sensitive information between developers.
 It is also how `exo deploy` interpolates secret variables into Terraform files during deployment (see (documentation/commands/deploy.md)).

#### use with `exo deploy`
 To create a secret to be used in deployment, first define its name under `remote.environments.#{remote-environment-id}.secrets` of a service's `service.yml` file:

```
remote:
  environments:
    <remote-environment-id>:
      secrets:
        - MONGODB_USER
        - MONGODB_PW
```

Second, manage secrets using the subcommands described above. Secrets are stored on S3 in the format: `secret_key = secret_value`.
 `secret_key` must match the corresponding secret name listed under the corresponding `secrets` block in `service.yml`. The value of `secret_value` is
 the string that Terraform injects into the service during deployment.
