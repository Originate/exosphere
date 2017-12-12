# RDS dependency

_Production plugin for RDS-type dependencies_

`type: rds`

#### Template configuration in application:
The following fields should be populated by the user:
- `engine`: RDS engine of choice
- `engine-version`: Engine version number
- `allocated-storage`: Allocated storage in gigabytes
- `instance-class`: Instance type of RDS instance
- `db-name`: Name of database instance
- `username`: Username for master db user
- `storage-type`: Storage type, i.e. general purpose SSD, provisioned IOPS, magnetic, etc.
- `password-secret-name`: Name of secret mapped to the db password. See [exo configure documentation](https://github.com/Originate/exosphere/blob/master/documentation/commands/configure.md) for more on secrets.
- `service-env-var-db-name`: Env var key name to be used for the db name
- `service-env-var-username`: Env var key name to be used for the db username
- `service-env-var-db-password`: Env var key name to be used for the db password

Example:
```yml
# application.yml
remote:
  postgres: # db-key-name
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
      service-env-var-db-name: DB_NAME
      service-env-var-username: DB_USERNAME
      service-env-var-password: DB_PASSWORD
```


#### Service-specific dependency data:
There is no service-specific dependency data for RDS-type depenencies.
