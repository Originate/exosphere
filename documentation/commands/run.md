# exo run

_Runs an Exosphere application on the local machine_

Usage: `exo run`

- dockerizes all services and their dependencies (databases),
  so no installation of programming languages or runtimes is necessary.
- prepares them (installing dependencies, compiling)
- boots the dependencies up first, so that services see all dependencies running
- monitors for file system changes and reboots the affected services

An Exosphere application is run by running all of its services.
How each service is run is defined in the respective [service configuration]().

`Exo run` is built on top of [Docker Compose](https://docs.docker.com/compose).
