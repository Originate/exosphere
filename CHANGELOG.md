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
