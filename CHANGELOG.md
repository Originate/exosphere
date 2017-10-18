## 0.28.2 (2017-10-18)

#### Bug fixes

* `exo deploy`: tag images with repository URI and version before pushing to ECR

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
