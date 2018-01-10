# exo run

_Runs an Exosphere application on the local machine_

Usage: `exo run [flags]`

- dockerizes all services and their dependencies,
  so no installation of programming languages or runtimes is necessary.
- prepares them (installing dependencies, compiling)
- reboots services when they go down
- `exo generate docker-compose` is run at the beginning of this command

Flags:
- `--production`   Runs application in production mode by building `Dockerfile.prod` files. Uses startup commands listed in `Dockerfile.prod` instead of those in `service.yml`.

An Exosphere application is run by running all of its services.
Exosphere reads [aplication configuration](documentation/configuration_files/application.md) and [service configuration](documentation/configuration_files/service.md) to determine how to run an application.

`exo run` is built on top of [Docker Compose](https://docs.docker.com/compose).

