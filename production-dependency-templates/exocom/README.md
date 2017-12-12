# Exocom dependency

_Production plugin for RDS-type dependencies_

## Usage
- List database `exocom` as a remote dependency in `application.yml`, with `type: exocom`
  - Populate required `template-config` fields (defined below)
- List any exocom-specific service data under `dependency-data.exocom` in a service's `service.yml`

### Required fields
The following fields should be populated by the user.

#### Template configuration in application:
- `version`: Define which version of Exocom to use
Example:
```
# application.yml
remote:
  exocom:
    type: exocom
    template-config:
      version: 0.27.0
```

#### Service-specific dependency data:
List message configuration in each service's `service.yml` file

Example:
```
dependency-data:
  exocom:
    receives:
      - service.ping
    sends:
      - service.pong
```

### Generated fields
The following fields will be rendered automatically by Exosphere.

#### Template configuration in exocom/config.yml:
- `EXOCOM_HOST`: the endpoint which services use to reach Exocom. Should be in the format `exocom.#{application-name}.local`
- `docker-image`: specifies which Exocom docker image to use.

Example:
```
# config.yml
service-environment-variables:
  - EXOCOM_HOST: exocom.my-app.local

docker-image: origiante/exocom:0.27.0
```

