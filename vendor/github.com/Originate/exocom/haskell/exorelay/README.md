# exorelay-hs
Communication relay between Haskell code bases and the Exosphere environment


## How to Use

First, you will will need to import the module: `import Network.Exocom`

### Initialization
Then, initialize the system by calling:
```haskell
newExoRelay portNum serviceName listenPort errorHandler
```
where `portNum` is an int that represents the port that the exocom service is (already) listening on and `serviceName` is a string which represents the name of your service. `listenPort` is the port you want your service to listen on. `errorHandler` is a function of type `String -> IO ()`. Whenever an error occurs the string describing the error will call `errorHandler` asynchronously.
Be warned that the handler is called in a separate thread and therefore be careful about using non thread-safe functions and data.
Returned from this function is an `IO exoRelay` object which will be needed for future sending and receiving calls to the message bus.

### Sending Messages
To send a message, encode the message and command (message name) as a Data.Aeson.Value (Json version of your payload) and call:
```haskell
sendMsg exo cmd payload
```
where `exo` is the exorelay instance, `cmd` is the command, and `payload` is the message contents. For many datatypes, Data.Aeson (from the aeson package) can derive an instance of
toJSON with mimimal effort. Please see the Aeson library documentation for more details on
how to perform the conversion of a given datatype to JSON.

If you expect a reply to your sent message then call the following function instead:
```haskell
sendMsgWithReply exo cmd payload handler
```
where `exo`, `cmd`, and `payload` is the same as the previous function. `handler` is a function with a signature: `Value -> IO ()` where the input argument is the content that is received as a reply to your sent message. The input is the parsed JSON and from there you can marshall it to your own datatype using `decode` from the Data.Aeson library. Therefore, when a reply is received, your handler is called asynchronously on the given payload. Note that your handler is called in a different thread and therefore do not perform any thread-unsafe actions without proper precautions.

### Listening to Messages
To listen for messages, simply call:
 ```haskell
registerHandler exo cmd handler
```
where `exo` is the exorelay instance, `cmd` is the command to listen for. `handler` is a function of type `Value -> IO ()` which executes asynchronously when the listening system receives a message which matches the given command name. The passed in Value parameter is the payload of the received message in parsed json form. Essentially, if the listening subsystem sees a message with name parameter matching `cmd` it will execute the handler on the payload of that message. Again, since the handler is executed in a separate thread please don't perform any thread-unsafe operations without proper locking.

If you would like to listen for a message and then be able to send a reply please call this function:
```haskell
registerHandlerWithReply exo cmd handler
```
which is very similar to `registerHandler` except that the type of `handler` is `Value -> IO (String, Value)` which again executes if a message is received with name parameter matching the `cmd` of the function call. The handler then returns an IO tuple of the form `(retCmd, retPayload)` then, the system will send a reply with command (name) `retCmd` and payload `retPayload` to the sender of the message. Again, beware of performing non-thread-safe operations in the handler

### Errors
All errors are transmitted to the error handler function that is provided on construction of the exorelay object. Again, this handler is called from a separate thread so don't do any non-thread-safe activities.

# Caveats
Although the exorelay library itself is thread-safe and you may call any of its functions from any thread, many of the handlers are executed in separate threads therefore, beware of mutating any state of your application or calling non thread safe functions in your handlers without employing proper locking or knowing what you are doing.

# Building
* You will need zmq installed (at least version 4) as well as having access to an exocom instance.
* Then, just run clone the repo and run cabal build (hackage package coming soon)

# Testing
* Please clone the repo first
* Then, `cd` into the `test` directory
* Make sure that ports `4100` and `4001` are both open (if you need to change the ports please modify `test/index.coffee` and `test/Test.hs` with the new ports)
* Run `npm i` to install the mock exocom service
* run `npm start` to start the mock exocom service
* Then, in a separate terminal (or at least while the mock exocom service is still running), navigate to the top level of the repo and run `cabal test`

# Contributing
* All contributions are welcome
* Please fill out issues as you find them
* PRs are welcome and should be compared with the `master` branch
