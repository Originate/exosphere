# Exocom dependency

_Production plugin for the Exocom dependency_

`type: exocom`

#### Template configuration in application:
The following fields should be populated by the user:

Template configuration:
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
  environment:
    EXOCOM_HOST: exocom.<app-name>.local
```

Environment variables:
- `EXOCOM_HOST`: endpoint at which Exocom can be reached. Must be set to `exocom.<app-name>.local` as a global env var

Example:
```yml
# application.yml
remote:
  environment:
    EXOCOM_HOST: exocom.<app-name>.local
```

#### Service-specific dependency data:
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
