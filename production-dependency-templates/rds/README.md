# Exocom dependency

_Production plugin for RDS-type dependencies_

## Usage
- List database `#{db-key-name}` as a remote dependency in `application.yml`, with `type: rds`
  - Populate required `template-config` fields (defined below)
- Populate `config.yml` with the required fields (defined below)
- List any db-specific service data under `dependency-data.#{db-key-name}` in a service's `service.yml`

### Required fields
The following fields should be populated by the user.

#### Template configuration in application:
- `engine`: RDS engine of choice
- `engine-version`: Engine version number
- `allocated-storage`: Allocated storage in gigabytes
- `instance-class`: Instance type of RDS instance
- `db-name`: Name of database instance
- `username`: Username for master db user
- `storage-type`: Storage type, i.e. general purpose SSD, provisioned IOPS, magnetic, etc.
- `password-secret-name`: Name of secret mapped to the db password. See [exo configure documentation](https://github.com/Originate/exosphere/blob/master/documentation/commands/configure.md) for more on secrets.

Example:
```
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


### Generated fields
The following fields will be rendered automatically by Exosphere.

#### Template configuration in rds/config.yml:
- `#{db-key-name}_HOST`: the endpoint which services use to reach this RDS instance. `db-key-name` will be converted to screaming snake case and preprended to `_HOST`.
- `docker-image`: Exocom docker image to use

Example:
```
# config.yml
service-environment-variables:
  - POSTGRES_HOST
  - DATABASE_NAME
  - DATABASE_USERNAME
  - DATABASE_PASSWORD
```
