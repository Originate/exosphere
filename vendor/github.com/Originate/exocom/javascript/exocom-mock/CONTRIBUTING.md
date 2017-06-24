# ExoComMock-JS Developer Guidelines

## Install

* `npm i`
* add `./bin/` and `./node_modules/.bin/` to your PATH


## Development

You can run `watch` in a separate terminal to auto-compile changes.
That's not necessary when running tests via `spec`


## Testing

* run all tests: `spec`


## Update Dependencies

* `update-check`: checks whether dependencies need to be updated
* `update`: updates the dependencies to the latest version


## Deploy a new version

```
publish <patch|minor|major>
```
