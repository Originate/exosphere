#Mock implementation of ExoCom in Go
A mock implementation of [ExoCom-Dev](https://github.com/Originate/exocom-dev) for sending and receiving messages to your ExoServices in a test environment.

###Installation

```
$ go get github.com/Originate/exocomMock
```

Then simply import it into your testing environment with

```go
import "github.com/Originate/exocomMock"
```

###Usage

```go
exocom = exocomMock.New()
go exocom.Listen(4100)
```
