# go-execplus

An abstraction around [os/exec.Cmd](https://golang.org/pkg/os/exec/#Cmd)
that allows you to:

* wait for specific text to appear in the output
* receive output chunks via a channel

See the tests for examples
