# service.yml

```yml
# `private` or `public` or `worker` 
type:

# configuration for Dockerfile.dev
development:
  port: # container port to expose
  scripts:
    run: # script executed in `exo run`
    test: # script executed in `exo test`

# configuration for Dockerfile.prod
production:
  port: # container port to expose
  health-check: # the health check path (`public` type only)


# configuration for locally running the service
local: # same as `local` in `application.yml` but applies only to this service

# configuration for deployments of the service
remote:
  cpu: # see notes below
  memory: # see notes below
  dependencies: # same as `remote.dependencies` in `application.yml` but applies only to this service
  environments:
    <id>: # match the remote environment id in `application.yml`
      url: # url for the service (`public` type only)
      environment-variables: # map of environment variables to pass to this service
      secrets: # array of secret keys variables to pass to this service (see `exo configure`)

# map of data to pass to particular dependencies
dependency-data:
  <dependency-id>: # match the dependency id from `application.yml` or `service.yml`
    # can have any structure that can be converted into JSON
```

#### Notes:

* `remote.cpu`: Number of CPU units to reserve for service container (see "cpu" under [Container Definitions/Environment](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#container_definition_environment))
* `remote.memeory`: The hard limit (in MiB) of memory to allocate to the service container (see "memory" under [Container Definitions](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#container_definitions))

# Service types

## Private
Private services expose an internal http endpoint. Container ports must be listed in `development.port` and `production.port`.

Endpoints available in local development:
  - `#{SERVICE_ROLE}_INTERNAL_ORIGIN`: The internal http endpoint at which a service can be reached. Used for internal communication with other services.

Endpoints available in deployment:
  - `#{SERVICE_ROLE}_INTERNAL_ORIGIN`: The internal http endpoint at which a service can be reached. Terraform manages creation of these records in route53. Used for internal communication with other services.

## Public
Public services expose an external and internal http endpoint. Container ports must be listed in `development.port` and `production.port`.

Endpoints available in local development:
  - `#{SERVICE_ROLE}_EXTERNAL_ORIGIN`: The external http endpoint at which a public service can be reached. Exosphere automatically picks an available host port and binds it to the specified service container port.
  - `#{SERVICE_ROLE}_INTERNAL_ORIGIN`: The internal http endpoint at which a service can be reached. Used for internal communication with other services.

Endpoints available in deployment:
  - `#{SERVICE_ROLE}_EXTERNAL_ORIGIN`: The load-balanced external https endpoint, as defined by the URL listed in the `remote.envorinments.#{remote-environment-id}` block of `service.yml`. Terraform manages creation of these records in route53.
  - `#{SERVICE_ROLE}_INTERNAL_ORIGIN`: The internal http endpoint at which a service can be reached. Terraform manages creation of these records in route53. Used for internal communication with other services.

## Worker
Worker services have the option of exposing an internal non-http, non-load-balanced endpoint. If a worker service exposes a port, the container port must be listed in `development.port` and `production.port`.

Endpoints available in local development:
  - `#{SERVICE_ROLE}_HOST`: The internal endpoint at which a service can be reached

Endpoints available in deployment:
  - `#{SERVICE_ROLE}_HOST`: The internal endpoint at which a service can be reached


#### Notes for deployments

For the worker to be reached reliably in deployments, please use of the `update-route53` binary available [here](https://github.com/Originate/exosphere/releases). The binary should be copied into the worker service's Docker image and ran as part of that service's start up script. Here is an example start up script:
```bash
#!/usr/bin/env bash
set -e

if [ -n "$INTERNAL_HOSTED_ZONE_NAME" ]; then
  update-route53 $ROLE $INTERNAL_HOSTED_ZONE_NAME
fi
node src/server.js
```
* The `if` statement is used to prevent any attempt to update Route 53 when running the application locally
* `$ROLE` is automatically passed to every service
* `$INTERNAL_HOSTED_ZONE_NAME` must be set by the user as an environment variable in each remote environment block, and they must be of the form `#{service-role}.#{remote-environment-id}-#{app-name}.local`
