# Exosphere configuration files

Exosphere requires configuration files in the code base
in order to know about your application.
There are two configuration files:

- The application itself is configured via the [application.yml](application_yml.md) file.
  It lists the different code bases that this application contains.
- Each service is configured by a [service.yml](service_yml.md) file.
  It defines where the source code of the service is,
  how to start or test it, etc.
