# Exoservice SDK

## Install

* clone this repo
* `setup`
* add `./bin/` and `./node_modules/.bin` to your PATH


## Development

* the CLI runs against the compiled JS, not the source LS,
  so run `watch` in a separate terminal to auto-compile changes


## Testing

* run all tests: `spec`
* run unit tests: `tests`
* run linters: `lint`


## Update dependencies

```
$ update
```


## Deploy a new version

```
$ publish <patch|minor|major>
```
