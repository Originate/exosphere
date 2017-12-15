# RDS dependency

_Production plugin for RDS-type dependencies_

`type: rds`

#### Template configuration in application:
The following fields should be populated by the user:

Template configuration:
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

Environment variables:
- `DB_NAME`: environment variable for db name, must be the same as the one listed in `<dependency-id>.db-name`
- `DB_USERNAME` environment variable for db username, must be the same as the one listed in `<dependency-id>.username`

Example:
```yml
# application.yml
remote:
  environment:
    DB_NAME: my-db # must be the same as `<dependency-id>.db-name`
    DB_USERNAME: originate-user # must be the same as `<dependency-id>.username`
```

Secrets:
- Database password environment variable name. Must be the same as the one listed in `<dependency-id>.password-secret-name`

Example:
```yml
# application.yml
remote:
  secrets:
    - POSTGRES_PASSWORD # must be the same as `<dependency-id>.password-secret-name`
```
