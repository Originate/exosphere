# Exoservice SDK

## Install

* clone this repo
* `npm i`
* add `./bin/` to your PATH


## Development

* the CLI runs against the compiled JS, not the source LS,
  so run `watch` in a separate terminal to auto-compile changes


## Testing

* run all tests: `spec`
* run unit tests: `tests`
* start watcher: `watch`
* run linters: `lint`


## Update dependencies

```
david
```


## Deploy a new version

```
npm version <patch|minor|major>
npm publish
```
