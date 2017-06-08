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
    ```

#### New Features

* New command `exo-clean`: removes dangling Docker images and volumes
* `exo-add`: add go template
