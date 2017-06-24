# Exoservice.js Developer Guidelines

## Install

* `npm i`
* add `./bin/` to your PATH


## Development

* the CLI runs against the Webpack build of exoservice, not the source LS,
  so run `webpack --watch` in a separate terminal to auto-compile changes


## Testing

```
$ spec
$ lint
```

- set environment variable `DEBUG_EXOSERVICE=1` to print debug output

## Update

```
$ update
```


## Deploy a new version

```
$ publish <patch|minor|major>
```
