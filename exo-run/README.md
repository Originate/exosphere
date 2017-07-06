# "exo run" command

This repo contains the code for the "exo run" command.
This is part of the [Exosphere SDK](https://github.com/Originate/exosphere-sdk)

The command will run the application on your machine.
  Exosphere boots up all services and their dependencies
  within Docker images, so no installation
  of programming languages or runtimes is necessary.
  It is built on top of [Docker Compose](https://docs.docker.com/compose/)
  and adds missing features on top of it.
  For example waiting until dependencies are fully booted up
  before starting the application services
  or restarting individual Docker images
  as their content gets modified by the developer.

More information available in [features](features).

Runs on top of public and private cloud environments
  like [AWS](https://aws.amazon.com),
  [Google Cloud](https://cloud.google.com),
  [Azure](https://azure.microsoft.com),
  [DCOS](https://dcos.io), etc<sup>&#42;</sup>

<hr>

<sup>&#42;</sup>
coming soon
