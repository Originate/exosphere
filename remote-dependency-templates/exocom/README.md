# Exocom dependency

_Production plugin for the Exocom dependency_

`type: exocom`

## Template configuration
The following fields are reuqired
- `version`: Define which version of Exocom to use

Example:
```yml
# application.yml
remote:
  dependencies:
    <dependency-id>:
      type: exocom
      template-config:
        version: 0.27.0
```

## Environment variables
In order for services to connect to exocom, `EXOCOM_HOST` should be set as an environment variable

Example:
```yml
# application.yml
remote:
  environments:
    <remote-environment-id>:
      EXOCOM_HOST: exocom.<remote-environment-id>-<app-name>.local
```

## Service data
List message configuration in each service's `service.yml` file and add translations in `application.yml`

Example:
```yml
# application.yml
services:
  service-a:
    dependency-data: # optional message translation block
      <dependency-id>:
        translations:
          - public: service-a ping
          - internal: internal-service pong
```
```yml
# service.yml
dependency-data:
  <dependency-id>:
    receives:
      - service.ping
    sends:
      - service.pong
```
