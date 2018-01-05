# application.yml

```yml
# application name - only lowercase alphanumeric character(s) separated by a single hyphen are allowed
name:

# configuration for locally running the application
local:
  dependencies:
    # each dependency has its own id
    <id>:
      # docker image (preferrably with tag)
      image:
      # map of environment variables to pass to this dependency
      environment-variables:
      # array of shell environment variables to pass through to this dependency
      secrets:
      # array of paths within the docker images you would like saved in a docker volume
      persist:

  # map of environment variables to pass to every service
  environment-variables:

  # array of shell environment variables to pass to all services
  secrets:

# configuration for deployments of the application
remote:
  dependencies:
    # each dependency has its own id
    <id>:
      # dependency type (exocom or rds)
      type:
      # map of data to pass to the dependency template
      template-config:
  environments:
    # each remote environment has its own id
    <id>:
      # map of environment variables to pass to every service
      environment-variables:
      # array of secret keys variables to pass to every service
      # use `exo configure` to edit
      secrets:
      # url for the application
      url:
      # aws account id
      account-id:
      # aws region
      region:
      # aws ssl cetificate arn
      ssl-certificate-arn:

# map of services, used locally and in deployments
services:
  # each service has its own role
  <role>:
    # relative path to the service source (use this or `docker-image`, not both)
    location:
    # docker image of the service source (use this or `location`, not both)
    docker-image:
```
