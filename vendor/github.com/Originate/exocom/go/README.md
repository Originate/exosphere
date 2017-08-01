# Development

The repo needs to live in a particular place in order for go to find everything properly

* set the environment variable `$GOPATH` to your Go workspace (typically `~/go`)
* `mkdir -p $GOPATH/src/github.com/Originate`
* `cd $GOPATH/src/github.com/Originate`
* `git clone git@github.com:Originate/exocom.git`

Now you should be able to run into any of the package directories and run the tests.

See [here](https://golang.org/doc/code.html) for more documentation
