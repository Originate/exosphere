# RDS dependency

_Production plugin for RDS-type dependencies_

`type: rds`

## Template configuration
The following fields are required:
- `engine`: RDS engine of choice
- `engine-version`: Engine version number
- `allocated-storage`: Allocated storage in gigabytes
- `instance-class`: Instance type of RDS instance
- `db-name`: Name of database instance
- `username`: Username for master db user
- `storage-type`: Storage type, i.e. general purpose SSD, provisioned IOPS, magnetic, etc.
- `password-secret-name`: Name of secret mapped to the db password. See [exo configure documentation](https://github.com/Originate/exosphere/blob/master/documentation/commands/configure.md) for more on secrets.

Example:
```yml
# application.yml
remote:
  dependencies:
    <dependency-id>:
      type: rds
      template-config:
        engine: postgres
        engine-version: 0.0.1
        allocated-storage: 10
        instance-class: db.t2.micro
        db-name: my-db
        username: originate-user
        storage-type: gp2
        password-secret-name: POSTGRES_PASSWORD
```

## Environment variables
The database name and username should be remote environment variables so that the RDS instance can be accessed. Name them as you choose but ensure their values match those defined in the template config. In order for services to connect to the instance, `<dependency-id>_HOST` should also be an environment variable.

Example:
```yml
# application.yml
remote:
  environments:
    <remote-environment-id>:
      DB_NAME: my-db # must match `db-name`
      DB_USERNAME: originate-user # must match `username`
      <dependency-id>_HOST: <db-name>.<remote-environment-id>-<app-name>.local
```

## Secrets
The database password should be a remote secret so that the RDS instance can be accessed. It should match the value of `password-secret-name`.

Example:
```yml
# application.yml
remote:
  secrets:
    - POSTGRES_PASSWORD # must match `password-secret-name`
```
