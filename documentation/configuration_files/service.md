# service.yml

```yml
# `worker` or `public`
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
