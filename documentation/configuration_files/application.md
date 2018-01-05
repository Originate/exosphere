# application.yml

```yml
# application name - only lowercase alphanumeric character(s) separated by a single hyphen are allowed
name:

# configuration for locally running the application
local:
  dependencies:
    <id>: # each dependency needs an id
      image: # docker image
      environment-variables: # map of environment variables to pass to this dependency
      secrets: # array of shell environment variables to pass through to this dependency
      persist: # array of paths within the docker images you would like saved in a docker volume
  environment-variables: # map of environment variables to pass to every service  
  secrets: # array of shell environment variables to pass to all services

# configuration for deployments of the application
remote:
  dependencies:
    <id>: # each dependency needs an id
      type: # dependency type (exocom or rds)
      template-config: # map of data to pass to the dependency template, see notes
  environments:
    <id>: # each environment needs an id
      account-id: # aws account id
      region: # aws region
      url: # url for the application
      ssl-certificate-arn: # aws ssl cetificate arn
      environment-variables: # map of environment variables to pass to every service
      secrets: # array of secret keys variables to pass to every service (see `exo configure`)

# map of services, used locally and in deployments
services:
  <role>: # each service has its own role
    location: # relative path to the service source (use this or `docker-image`)
    docker-image: # docker image of the service source (use this or `location`)
```

##### Notes

* `remote.dependencies.<id>.template-config`: see the `README.md` within each folder in [remote-dependency-templates](../../remote-dependency-templates) for more information
