# Scaffolding

## Creating a new application

```
exo create app my-app
```

* the script asks for all important information
* this creates the following file structure:

  ```
  config.yml           # application configuration file
  web/                 # put your web server stack here
  api/                 # put your API (e.g. GraphQL) server stack here
  services/            # a directory to contain application-specific microservices
  spec/features/       # end-to-end feature tests
  spec/performance/    # performance tests
  documentation/       # documentation for the app goes here, in MarkDown
  ```


## Creating a new service

```
exo create service my-service
```

* asks for extra information
* creates the following file structure

  ```
  services
