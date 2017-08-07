# Release Process

* Create a feature branch which updates
  * `CHANGELOG.md`
  * the version in `exo-go/src/cmd/version.go` and the related feature
* Get the feature branch reviewed and merged
* Create and push a new Git tag for the release
  * `git tag vX.Y.Z`
  * `git push --tags`
* Travis-CI creates a new release on Github and attaches the binaries to it
