# exo generate

_Generates docker compose and terraform files_

Usage: `exo generate [command]`

Flags:
- `-p, --profile string`  AWS profile to use (defaults to "default"). Only available for `exo generate terraform-var-files`

Available Commands:
-  `exo generate docker-compose`      Generates docker-compose files
-  `exo generate terraform`           Generates terraform files
-  `exo generate terraform-var-files` Generates terraform tfvar files

Flags:
- `exo generate docker-compose --check`   Runs check to see if docker-compose are up-to-date

Notes:
- `exo generate docker-compose` is automatically run before `exo run`, `exo test`, and `exo clean`
- `exo deploy` will fail unless `exo generate terraform` is run beforehand
