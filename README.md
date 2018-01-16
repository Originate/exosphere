<img src="documentation/logo.png" width="862" height="111" alt="logo">

[![Build Status](https://travis-ci.org/Originate/exosphere.svg?branch=master)](https://travis-ci.org/Originate/exosphere)

_yak shaver for cloud developers_

Exosphere helps work on cloud applications.
It automates the repetitive activities
to get them installed, built, running, tested, and deployed:

### Setting up an application
- [exo init](documentation/commands/init.md)
  initializes a new exosphere application
- [exo add](documentation/commands/add.md)
  adds a new service to the application

### Developing an application
- [exo generate docker-compose](documentation/commands/generate/docker-compose.md)
  generates docker compose files
- [exo run](documentation/commands/run.md)
  runs an Exosphere application on the local machine
- [exo test](documentation/commands/test.md)
  runs all the tests for an application
- [exo clean](documentation/commands/clean.md)
  cleans docker workspace

### Deploying an application
- [exo generate terraform](documentation/commands/generate/terraform.md)
  generates terraform file
- [exo generate terraform-var-file](documentation/commands/generate/terraform-var-file.md)
  generates terraform var file
- [exo configure](documentation/commands/configure.md)
  manages remotely stored application secrets
- [exo deploy](documentation/commands/deploy.md)
  pushes the application to public or private clouds

### Configuring Exosphere
- [exo template](documentation/commands/template.md)
  manages templates for services


#### Learn more
* [why Exosphere](documentation/benefits.md)
* which [open source](documentation/open-source.md) technologies Exosphere is built on top of
* the [configuration files](website/config_files)


#### Get started
* [install](website/tutorial/part_1/03_installation.md) it
* [download and run an example application](website/example-apps.md)
* build your own application by following the [tutorial](website/tutorial)


#### Get involved
* [platform developer documentation](website/developers/developers.md)
