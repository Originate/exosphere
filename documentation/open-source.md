# Built on top of open source

Exosphere is built on top of industrial-strength open-source technologies
and makes them work well together based on shared configuration:
* [Git](https://git-scm.com) for source code management
* [Docker](https://www.docker.com) for containerization.
  Exosphere automatically dockerizes all code bases and runs them inside Docker.
* [Docker Compose](https://docs.docker.com/compose) for booting up all application services at once.
  Exosphere adds a few missing features on top of Docker Compose,
  like improved startup orders and restarting services on filesystem changes.
* [Boilr](https://github.com/tmrts/boilr) for code generation of new services
* [Terraform](https://www.terraform.io) for infrastructure-as-code based deployment.
  Exosphere generates Terraform code based on application information,
  gives the users a chance to fine-tune it,
  and guides them through the process of standing up
  development, staging, production, and other environments using it.
