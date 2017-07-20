# Service Configuration

Each Exosphere service is configured via a `service.yml` file.
This file exists in the root directory of each service.
It provides information about the service,
and can contain the following elements:

- __type:__ the type of service.
  This is comparable to the class name in object-orientation.
- __description:__ textual description for this service
- __author:__ who wrote this service (name and email)
- __setup:__ command to prepare the code base of the service after cloning it from its repository,
  for example `npm install` for Node.js or `go install` for Golang code bases.
- __startup:__ This section defines how to start the service.
  This is done via two elements:
  - __command:__ the command to start the service up inside its Docker image.
    For example `npm run`
  - __online-text:__ the console output to wait for
    to know that the service is fully booted up
    and ready to accept connections
- __docker:__ Things you want to add to this service's _Dockerfile_,
  for example exposing ports.

Next in the file comes hosting-specific configuration data.
For example, when hosting on AWS, you need to provide:
- __aws:__
  - __cpu:__ needed CPU ??? (default: 100)
  - __memory:__ needed memory ??? (default: 500)
  - __essential:__ ??? (default is "true")
  - __public-port:__ for public services only.
    Specifies ???

If this service is public, it needs additional configuration for each environment (staging, production, etc):
- __url:__ the address under which the service is visible on the internet

For an example file, please check out the
[SpaceTweet web server configuration file](https://github.com/Originate/space-tweet/blob/master/web-server/service.yml).
