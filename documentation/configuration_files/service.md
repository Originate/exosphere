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
  health-check: # for `public` services the health check path


# configuration for locally running the service
local: # same as `local` in `application.yml` but applies only to this service

# configuration for deployments of the service
remote:
  cpu:
  memory:
  dependencies: # same as `remote.dependencies` in `application.yml` but applies only to this service
  environments:
    <id>: # match the remote environment id in `application.yml`
      url: # url for the application
      environment-variables: # map of environment variables to pass to this service
      secrets: # array of secret keys variables to pass to this service (see `exo configure`)

# map of data to pass to particular dependencies
dependency-data:
  <dependency-id>: # match the dependency id from `application.yml` or `service.yml`
    # can have any structure that can be converted into JSON
```
