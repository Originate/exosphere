# exo generate

_Generates docker compose and terraform files_

Usage: `exo generate [command]`

- `exo generate docker-compose` is automatically run before `exo run` and `exo test`
- `exo deploy` will fail unless `exo generate terraform` is run beforehand

Available Commands:
-  `exo generate docker-compose` Generates docker-compose files
-  `exo generate terraform`      Generates terraform files

Flags:
- `exo generate docker-compose --check`   Runs check to see if docker-compose are up-to-date


