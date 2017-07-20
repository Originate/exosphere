# Application Configuration

Each Exosphere application is configured via an `application.yml` file.
This file exists in the root directory of the application.
It provides information about the application itself,
and can contain the following elements:

- __name:__ the application name
- __description:__ a description of the application
- __version:__ the version of the application

- __dependencies:__ this block specifies application-wide dependencies,
  for example the message bus or discovery service used to enable inter-service communication.
  Exosphere must know about the dependency types,
  for example what type of information to tell the services about it.

  Dependencies are specified via two keys:
  - __type:__ the dependency type
  - __version:__ which version of the dependency to use

  Exosphere supports the following dependencies:
  - __exocom:__ Originate's micro-service optimized message bus
  - __zookeeper:__ Apache Zookeeper (coming soon)

- __services:__ defines the services that the application is comprised of.
  Exosphere supports two types of services:
  - __public:__ services that are exposed to the internet,
    for example web or API servers.
  - __private:__ services in a private network, only accessible by the public services,
    for example storage or compute services

  Each service entry is keyed with the _role_ that the service plays in this application,
  for example "web server", "users service", etc.
  Services are configured via these entries:
  - __location:__ where the code lives on the machine, relative to the application's root folder.
    (default: `./<service name>`
  - __docker_image:__ name of the docker image for this service's code base on DockerHub

- __environments:__ defines the different hosting environments for this application,
  for example `development`, `staging`, or `production`.
  Each environment can be configured via the following entries:
  - __url:__ the URL under which this environment is running,
    for example `staging.acme.com`
  - __provider:__ the cloud provider to deploy into.
    The following values are possible:
    - __aws:__ Amazon AWS. The AWS driver can be configured via following entries:
      - __region:__ which AWS region to deploy into
      - __public-cluster-size:__ machine size to be used for the public ECS cluster (default is `t2.nano`)
      - __private-cluster-size:__ machise size to be used for the private ECS cluster (default is `t2.nano`)
      - __remote-state-store:__ name of the S3 bucket in which to store the Terraform state (default is `terraform.state.<environment name>`)
    - __gc:__ Google Cloud (coming soon)
    - __azure:__ Azure
    - __dcos:__ a private DCOS cloud environment


For an example file, please check out the
[SpaceTweet configuration file](https://github.com/Originate/space-tweet/blob/master/application.yml).
