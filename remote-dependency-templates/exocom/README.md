# Exocom dependency

_Production plugin for the Exocom dependency_

`type: exocom`

#### Template configuration in application:
The following fields should be populated by the user:
- `version`: Define which version of Exocom to use

Example:
```yml
# application.yml
remote:
  <dependency-id>:
    type: exocom
    template-config:
      version: 0.27.0

services:
  service-a:
    dependency-data: # optional message translation block
      <dependency-id>:
        translations:
          - public: service-a ping
          - internal: internal-service pong
```

#### Service-specific dependency data:
List message configuration in each service's `service.yml` file

Example:
```yml
dependency-data:
  <dependency-id>:
    receives:
      - service.ping
    sends:
      - service.pong
```
