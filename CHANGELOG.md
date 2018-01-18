## Unreleased

## 0.36.1 (2017-01-17)

#### Bug fixes

* `exocom` dependency terraform: fix cluster - instance connection

## 0.36.0 (2017-01-15)

#### BREAKING CHANGES

* storage of secrets on s3 changed
  * `{{account-id}}-{{app-name}}-{{remote-environment-id}}-terraform-secrets/secrets.json` ->  `{{account-id}}-{{app-name}}-{{remote-environment-id}}-terraform/secrets.json`
* update infrastructure names in `exocom` dependency to include remote environment id
* update internal network name to include remote environment id

#### New Features

* add support for shared code between services
  ```yml
  # application.yml

  # array of directory names, each of which is copied into every service before building the docker image
  shared:
  ```
* worker services can expose ports. See [here](https://github.com/Originate/exosphere/blob/0721a4d527cb3d2bde864c52a6e5374fbe8299d9/documentation/configuration_files/service.md#worker) for documentation

#### Bug fixes

* `exo test`: stop printing `The following tests failed:` when no tests fail

## 0.35.0 (2017-12-22)

#### BREAKING CHANGES
* Rename `environment` to `environment-variables`
  ```yml
  # application.yml / service.yml
  local:
    environment:

  remote:
    environment:

  # becomes
  local:
    environment-variables:

  remote:
    environment-variables:
  ```
* Support deploying to multiple environments.
  * `exo configure` and `exo deploy` now require the remote environment id to be the first argument
  * storage of secrets and the terraform state on s3 changed
    * `{{account-id}}-{{app-name}}-terraform-secrets/secrets.json` ->  `{{account-id}}-{{app-name}}-{{remote-environment-id}}-terraform-secrets/secrets.json`
    * `{{account-id}}-{{app-name}}-terraform/terraform.tfstate` ->  `{{account-id}}-{{app-name}}-{{remote-environment-id}}-terraform/terraform.tfstate`
    * Also update the `LockID` in the DynamoDB TerraformLocks table to reflect the new path
  * configuration updates
  ```yml
  # application.yml
  remote:
    dependencies:
    url:
    region:
    account-id:
    ssl-certificate-arn:
    environment-variables:
    secrets:

  # becomes

  remote:
    dependencies:
    environments:
      <remote-environment-id>: # for example qa or production
        url:
        region:
        account-id:
        ssl-certificate-arn:
        environment-variables:
        secrets:
  ```
  ```yml
  # service.yml
  remote:
    dependencies:
    environment-variables:
    secrets:
    url:
    cpu:
    memory:

  # becomes

  remote:
    dependencies:
    cpu:
    memory:
    environments:
      <remote-environment-id>: # for example qa or production
        url:
        environment-variables:
        secrets:
  ```
* aws: update log bucket from `production-{{app-name}}-logs` to `{{account-id}}-{{app-name}}-{{remote-environment-id}}-logs`

#### New Features
* `exo deploy`: use terraform `-var-file` instead of `-var`

## 0.34.0 (2017-12-18)

Fixes versioning from previous release (should have been major release)

#### Bug Fixes
* Fix required fields for RDS dependency

## 0.33.3 (2017-12-18)

#### BREAKING CHANGES
* remote dependencies have been generalized, please refer to `remote-dependency-templates/#{dependency-type}/README.md` for dependency-specific details. Example change:
```yml
remote:
  dependencies:
    exocom:
      type: exocom
      config:
        version: 0.27.0

# becomes

remote:
  dependencies:
    exocom:
      type: exocom
      template-config:
        version: 0.27.0
```
* Remove manual port management for local dependencies:
```yml
mongo:
  image: mongo:3.4.0
  ports:
    - `4000:4000`

# becomes

mongo:
  image: mongo:3.4.0
```
* `exo-deploy` has been reverted to run Terraform on local machine


#### New Features
* application-wide environment variables and secrets now supported. Example:
```yml
remote:
  environment:
    key: value
  secrets:
    - secret-key
```

#### Bug Fixes
* throw previously ignored `exo run` errors
* generate Terraform dependency modules in deterministic order

## 0.33.2 (2017-12-12)

* Fix invalid release 0.33.1 (released the wrong code)

## 0.33.1 (2017-12-12)

#### Bug Fixes
* Fix passing of dependency data to remote a exocom dependency

#### New Features
* Type is no longer necessary for local dependencies
```yml
# Before
local:
  dependencies:
    <id>:
      type: # exocom, nats, or generic (defaults to generic if omitted)
      image:
      config:

# After
local:
  dependencies:
    <id>:
      image:
      config:
```

## 0.33.0 (2017-12-11)

#### BREAKING CHANGES
* Local and remote dependencies have been refactored to be more generic:
```yml
local:
  - name:
    version:
    config:

remote:
  - name:
    version:
    config:

# becomes

local:
  <id>: # user-defined dependency key
     type: # exocom, nats, or generic (defaults to generic if omitted)
     image: # always required
     config:

remote:
  <id>: # user-defined dependency key
     type: # exocom, rds, or generic (defaults to generic if omitted)
     config:
```

* Remove environemnts block:
```yml
environment:
  default:
  local:
  remote:
  secrets:

# becomes:

local:
  environment:

remote:
  environment:
  secrets:

# default is removed in favor of using yaml defaults
```

* Extract exocom messages into generalized `dependency-data` block:
```yml
# service.yml
messages:
  receives:
    - service.ping
  sends:
    - service.pong

# becomes:
dependency-data:
  exocom:
    receives:
      - service.ping
    sends:
      - service.pong
```
And
```yml
# application.yml
services:
  users-service:
    message-translation:
      - public: users ping
        internal: mongo ping

# becomes
services:
  users-service:
    dependency-data:
      exocom:
        translations:
          - public: users ping
            internal: mongo ping
```
* Move `health-check` field from remote into production block
```yml
remote:
  health-check:

# becomes
production:
  health-check:
```

#### New Features
* For each public service, inject `<PUBLIC_SERVICE>_INTERNAL_ORIGIN` environment variable to every other service. This points to exoposed origins on an internal network
* Switch local secrets to using variables in generated docker-compose files
* `exo generate`: check for application and service config schema errors
* Remote container names from docker-compose files
* Fix app config validation error messages
* Validate Terraform files before beginning to push Docker images to ECR

## 0.32.0 (2017-12-01)

#### BREAKING CHANGES
* `exo deploy`
  * rename flag `--update-services` to `--auto-approve`
  * updated to always errors if the terraform file is not up to date
* `exo run`: remove `--no-mount` flag, now always mounts
* `exo test`: remove `--no-mount` flag, now never mounts
* `exo create`: removed in favor of `exo init` which initializes the current directory as an exosphere project
* dependency volumes have been updated to use named docker volumes
  ```yml
  # Before
  - name: 'mongo'
    version: '3.4.0'
    config:
      volumes:
        - '{{EXO_DATA_PATH}}:/data/db'

  # After
  - name: 'mongo'
    version: '3.4.0'
    config:
      persist:
        - /data/db
  ```
* update service environment variables from development/production to local/remote
  ```yml
  # Before
  environment:
    development:
      ENV2: value2
      ENV3: dev_value3
    production:
      ENV3: prod_value3

  # After
  environment:
    local:
      ENV2: value2
      ENV3: dev_value3
    remote:
      ENV3: prod_value3
  ```
* break up application and service development/production blocks
  ```yaml
  ########################################
  # application.yml
  ########################################
  # Before
  development:
    dependencies:

  production:
    dependencies:
    url:
    region:
    account-id:
    ssl-certificate-arn:

  # After
  local:
    dependencies:

  remote:
    dependencies:
    url:
    region:
    account-id:
    ssl-certificate-arn:

  ########################################
  # service.yml
  ########################################
  # Before
  development:
    dependencies:
    port:
    scripts:

  production:
    dependencies:
    port:
    url:
    cpu:
    memory:
    health-check:

  # After
  development:
    port:
    scripts:

  local:
    dependencies:

  production:
    port:

  remote:
    dependencies:
    url:
    cpu:
    memory:
    health-check:
  ```
* terraform: inject dependency docker images as variables
* update to terraform 0.11.0 and run terraform in a docker container
* update cloudwatch alarm thresholds from 10/90 to 20/70

#### New Features
* add `exo generate` command
  * `exo generate docker-compose` generates the docker compose files
    * use the `--check` flag to verify the files are up to date
    * the files are also updated on each run of `exo run`, `exo clean`, `exo test`, and `exo deploy`
  * `exo generate terraform` generates the terraform files
* for each public service, all other services receive: `<SERVICE>_EXTERNAL_ORIGIN` as an environment variable which points to the exposed origin. `<SERVICE>` is the service role converted to constant case.
* update output of commands to include the directory its run in and the environment variables passed to it

## 0.31.0 (2017-11-13)


#### BREAKING CHANGES
* Change the way dependencies are handled. Dependencies and their configuration should be defined either in service.yml, if only that service uses it, or in application.yml, if more than one service uses it.
* Remove `exo template add`, `exo template fetch`, and `exo template remove`. `exo template test` remains supported.
* Move service template directory from `${APP_DIR}/.exosphere` to `${APP_DIR}/.exosphere/service_templates`
* Move exosphere data path from `${HOME_DIR}/.exosphere/...` to `${APP_PATH}/.exosphere/data`
* Service protection levels have been weaned down to `public` and `worker` only, and they should be defined in `service.yml` under `type`, rather than in `applicaiton.yml`

BEFORE:
```
# application.yml
services:
  public:
    service1:
      ...
  worker:
    service2:
      ...
```

AFTER:
```
# application.yml
services:
  service1:
   ...
  service2:
   ...
```
```
# service1 service.yml
type: public
```
```
# service2 service.yml
type: worker
```

#### New Features
* `exo-test`
  * run tests in deterministic order
  * print service names with failed tests
* generate `docker-compose` files in more deterministic manner
  * sort `depends_on` fields
  * remove user paths

#### Bug fixes
* `exo-run`: stop printing "docker-compose down" exit error when shutting down application with sigint
* `exo-deploy`: pass file name to docker-compose processes

## 0.30.0 (2017-11-07)

#### BREAKING CHANGES

* Assign public ports in development/production blocks of `service.yml`. `docker` block has been removed. Only a single port is supported at this time.
```yml
docker:
  ports:
    - '3000:3000'

production:
  public-port: 3000
```
Changes to:
```yml
development:
  port: 3000

production:
  port: 3000
```

#### Bug fixes

* `exo-test`: fix bug where every service was being built when testing only one service

## 0.29.0 (2017-11-06)

#### BREAKING CHANGES

* `exo-run`
  * removed ability to silence services and dependencies
  * remove concept of online texts
    * no longer spinning up dependencies first and waiting for them to come online before spinning up services
* `exo-run` and `exo-test`: generate named docker-compose files in `#{app-directory}/docker-compose` directory to be committed to git

#### New Features

* `exo-run`: restart services on failure
* `exo-deploy`
  * name ECS instance profiles and role services the same thing. Improves deployment teardown
  * print image push progress
* `exo-test`: catch sigint and gracefully shut down test containers
* `exo-clean`, `exo-deploy`, `exo-run`, and `exo-configure` commands can now be run in service directories
* improve output logs across all commands
  * print commands/timestamps of subprocesses
* build docker-compose network name from application name instead of application direcotory

## 0.28.4 (2017-10-25)

#### New Features

* `exo-deploy`: allow traffic from bastion instances to postgres RDS instances

## 0.28.3 (2017-10-19)

#### New Features

* `exo-run`: support secret env vars in development

## 0.28.2 (2017-10-18)

#### Bug fixes

* `exo deploy`: tag images with repository URI and version before pushing

## 0.28.1 (2017-10-18)

#### New Features
* `exo-deploy`
  * pass exocom service routes as a variable to main.tf file
  * don't build images that won't be pushed to ECR

#### Bug fixes

* `exo test`: fix bug that required service role and service directory be the same name
* `exo-deploy`: fix Terraform issue of ECS services being spun up before target groups were created

## 0.28.0 (2017-10-10)

#### BREAKING CHANGES

* exo-deploy: Rename S3 bucket and state files:
  * `#{app-name}-terraform` -> `#{account-id}-#{app-name}-terraform`
  * `#{app-name}-terraform-secrets` -> `#{account-id}-#{app-name}-terraform-secrets`
  * `dev/terraform.tfstate` -> `terraform.tfstate`

#### New Features

* add validation of application.yml fields
* exo-deploy: support loading aws credentials via environment variables

## 0.27.0 (2017-10-07)

#### BREAKING CHANGES

* Update rds config to auto inject variables to an services that dependent on it

#### New Features

* `exo template test`: use `--no-mount` when running `exo test`
* `exo deploy`
  * reduce number of things that cause changes in `main.tf`
    * pass service image as variable
    * put services in alphabetical order
    * pass aws profile as variable
  * revert changes to the `main.tf` file if deployment is abandoned
  * generate `terraform/.gitignore`
  * use commit sha as service version number
  * skip pushing images to ECR if version already exists
  * add flag `--update-services` which exits with an error if it causes any changes in `main.tf` and otherwise deploys updates without a confirmation (for use on CI)

#### Bug fixes

* Make external DNS a public hosted zone
* Compile service dependency env vars in deployment
* `exo clean`: allow killing containers to finish before exiting

## 0.26.3 (2017-09-26)

#### New Features

* `exo deploy`
 * validate required RDS fields for applicable dependencies
 * cleanup log messages

#### Bug fixes

* `exo deploy`
  * iron out remaining bugs in RDS module
  * ignore certain dependencies when pushing images to ECR
  * switch ECS containers to using to soft memory limits

## 0.26.2 (2017-09-26)

#### New Features

* `exo test`: fix env passed to docker-compose (caused failures on CI)

## 0.26.1 (2017-09-25)

#### New Features

* `exo run` / `exo test`: add `--no-mount` flag to disable mounting
* `exo deploy`
  * RDS support in production for postgres and mysql engines
  * Support for production dependencies listed only in `service.yml`

## 0.26.0 (2017-09-21)

#### BREAKING CHANGES

* Restructure configuration of dependencies in `application.yml` and `service.yml`
  ```yml
  dependencies:
  - name: ...
    version: ...

  # becomes

  development:
    dependencies:
    - name: ...
      version: ...

  production:
    dependencies:
    - name: ...
      version: ...
  ```

#### New Features

* `exo run`: add `--production` flag to use the production docker images

## 0.25.1 (2017-09-18)

#### Bug fixes

* `exo test`: fix for services that use mounting

#### Documentation

* add `exo configure` and `exo deploy` documentation

## 0.25.0 (2017-09-18)

#### BREAKING CHANGES

* `exo run`: replace restarting services with mounting
  * each service is mounted in `/mnt` on their dev docker container

## 0.24.0 (2017-09-15)

#### BREAKING CHANGES

* Service `Dockerfile` locations have changed
  * remove `Dockerfile` and `tests/Dockerfile`
  * add `Dockerfile.dev` for use in development (`exo run` and `exo test`)
  * add `Dockerfile.prod` for use in production

#### New Features

* `exo deploy`
  * validate production fields in `application.yml` and `service.yml`
  * read AWS region / account id / ssl certificate arn from `application.yml`

#### Bug fixes

* Service specific dependencies: support `config/service-environment`
* `exo configure`: can now handle spaces in secret values
* make prompts more clear
* `exo deploy`
  * better error message when app name conforms to s3 naming
  * use profile flag when generating terraform

## 0.23.0.alpha.5 (2017-08-30)

#### BREAKING CHANGES

* Drop support for `exocom@<0.26.1`

#### New Features

* add worker service
* `exo configure` / `exo deploy`: add profile flag to select AWS profile
* `exo deploy`
  * automate service environment variables
  * add ability to mark dependency as external in production
* `exo clean`: stop and remove all docker containers used for running / testing

#### Bug fixes

* `exo deploy`: remove requirement that service names be at least 32 characters
* `exo run` / `exo test`: use different docker networks to prevent clashes

## 0.23.0.alpha.4 (2017-08-30)

#### Bug fixes

* `exo run`: fix dependency environment variables

## 0.23.0.alpha.3 (2017-08-29)

#### New Features

* `exo deploy` ready for use
* `exo template test` command to test that a service template can be used and has passing tests

#### Bug fixes

* strip control characters from shutdown output

## 0.23.0.alpha.2 (2017-08-15)

#### BREAKING CHANGES

* Drop support for `exocom@<0.24.0`
* Update `exo clean` to use prune images and volumes docker commands
* Update `exo template fetch` to fetch updates for a single template that must be passed as a parameter

#### New Features

* Add `exo configure` to manage deployment secrets
* Add `exo test`
* Add support for `exocom@0.24.0`
* Update `exo template add` to take an optional commit / tag

#### Known issues

* incomplete `exo deploy` which is still being developed

## 0.23.0.alpha (2017-08-04)

#### BREAKING CHANGES

* Rewrite in Go, many commands added, changed, and removed.

#### Known issues

* incomplete `exo deploy` and missing `exo test` each of which are being developed

## 0.22.1 (2017-06-08)

#### Bug fixes
* fix [nats](http://nats.io/) application dependency

## 0.22.0 (2017-06-08)

#### BREAKING CHANGES

* Add support for application dependencies
  * If using exocom your `application.yml` requires the following update:
    ```yml
    bus:
      type: exocom
      version: 0.21.7

    # becomes

    dependencies:
      - type: exocom
        version: 0.21.7
    ```
  * Added support for using [nats](http://nats.io/) with the following in your `application.yml`
    ```yml
    dependencies:
      - type: nats
        version: 0.9.6
    ```
    The environment variable `NATS_HOST` will be passed to all services
* `exo-setup`  and `exo-run`: switch to [DockerCompose](https://docs.docker.com/compose/) from running individual `docker` processes
  * Requires the following updates to your `service.yml` files
    ```yml
    dependencies:
      mongo:
        version: '3.4.0'
        docker_flags:
          volume: '-v {{EXO_DATA_PATH}}:/data/db'
          online_text: 'waiting for connections'
          port: '-p 27017:27017'

    # becomes

    dependencies:
      mongo:
        dev:
          image: 'mongo'
          version: '3.4.0'
          volumes:
            - '{{EXO_DATA_PATH}}:/data/db'
          ports:
            - '27017:27017'
          online-text: 'waiting for connections'
    ```
    ```yml
    docker:
      publish:
        - '3000:3000'

    # becomes

    docker:
      ports:
        - '3000:3000'
    ```

#### New Features

* New command `exo-clean`: removes dangling Docker images and volumes
* `exo-add`:
  * add go template
  * prompt user for protection level of service
