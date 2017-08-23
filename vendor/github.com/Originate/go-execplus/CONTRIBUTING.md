# Release Process

* Create a feature branch with the following
  * Update `CHANGELOG.md`
  * Update version in `execplus.go`
* Get the feature branch reviewed
* Squash merge using the commit message `Release vX.Y.Z`
* Create and push a new Git tag for the release
  * `git tag vX.Y.Z`
  * `git push --tags`
