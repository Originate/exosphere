# Exoservice.js Developer Guidelines

## Install

* `npm i`
* add `./bin/` to your PATH


## Development

* the CLI runs against the compiled JS, not the source LS,
  so run `watch` in a separate terminal to auto-compile changes


## Testing

```
$ spec
$ lint
```

## Update dependencies

```
$ update
```


## Deploy a new version

```
$ publish <patch|minor|major>
```
