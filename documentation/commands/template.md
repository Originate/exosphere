# exo template

_Manages templates for Exosphere services_

Usage `exo template <[subcommand]>`

Available subcommands:
- `exo template add <name> <URL>`:
  downloads a [Boilr](https://github.com/tmrts/boilr) template
  and registers it under the given name.
  The template is added as a Git submodule to the application repository
- `exo template remove <name>`:
  removes the template with the given name from the application's template registry
- `exo template fetch`: fetches updates for all templates
