require! {
  'cucumber': {defineSupportCode}
  'prelude-ls' : {pairs-to-obj}
  'sinon'
  'wait' : {wait}
}


defineSupportCode ({When}) ->


  When /^a call comes in$/, (done) ->
    @create-websocket-endpoint {@exocom-port}, ~>
      @service-send-message {name: \foo, id: \123}
      done!


  When /^a new service instance registers itself with it via the message:$/ (table, done) ->
    table-data = table.raw! |> pairs-to-obj
    payload = table-data.PAYLOAD |> JSON.parse
    @create-named-websocket-endpoint {name: 'test instance', @exocom-port, registration-message: table-data.NAME, registration-payload: payload}, done


  When /^closing it$/, (done) ->
    @exocom.close done


  When /^I tell it to wait for a call$/, ->
    @call-received = sinon.spy!
    @exocom.on-receive @call-received


  When /^resetting the ExoComMock instance$/, ->
    @exocom.reset!


  When /^sending a "([^"]*)" message to the "([^"]*)" service with the payload:$/, (message, service, payload) ->
    @exocom-send-message {@exocom, service, message-data: {name: message, payload: payload}}


  When /^trying to send a "([^"]*)" message to the "([^"]*)" service$/, (message-name, service-name, done) ->
    try
      @exocom.send service: service-name, name: message-name
    catch
      @error = e
      done!
