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
      persist: # array of paths within the docker image you would like saved to docker volumes
  environment-variables: # map of environment variables to pass to every service
  secrets: # array of shell environment variables to pass to every services

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
      secrets: # array of secret keys used to lookup secrets to pass to every service (see `exo configure`)

# map of services, used locally and in deployments
services:
  <role>: # each service has its own role
    # Only one of location / docker-image should be specified, not both
    location: # relative path to the service source
    docker-image: # docker image of the service source

# array of directory names, each of which is copied into every service
shared:
```

##### Notes

* `remote.dependencies.<id>.template-config`: see the `README.md` within each folder in [remote-dependency-templates](../../remote-dependency-templates) for more information
