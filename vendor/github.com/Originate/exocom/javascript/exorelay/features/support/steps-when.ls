require! {
  'chai' : {expect}
  'ejs'
  'livescript'
  'sinon'
  'wait' : {wait, wait-until}
  '../..' : ExoRelay
  'exocom-mock': MockExoCom
  'prelude-ls' : {any}
  'cucumber': {When}
}


When /^an ExoCom instance comes online at port (\d+), (\d+) second(?:s)? later$/ timeout: 10_000, (port, seconds, done) ->
  wait seconds * 1000, ~>
    @exocom = new MockExoCom
      ..listen port
    done!


When /^ExoCom goes down for (\d+) second(?:s)? and comes back online$/ timeout:10_000, (seconds, done) ->
  @exocom.close ~>
    wait seconds * 1000, ~>
      @exocom = new MockExoCom
        ..listen @exocom-port, done


When /^an ExoRelay instance running inside the "([^"]*)" service comes online$/ (@role, done) ->
  @exo-relay = new ExoRelay {@role, @exocom-port, exocom-host: "localhost"}
    ..connect!
    ..on 'online', ~>
      @message-id = @exo-relay.websocket-connector.last-sent-id
      done!
    ..on 'error', (@error) ~>


When /^I check the status$/, (done) ->
  @exocom.send service: 'test-service', name: '__status'
  wait-until (~> @exocom.received-messages.length), 1, ~>
    @status-code = @exocom.received-messages[0].name
    done!


When /^creating an ExoRelay instance using ExoCom host "([^"]*)" and port (\d+)$/ (host, +port) ->
  @exo-relay = new ExoRelay do
    exocom-host: host
    exocom-port: port
    role: 'test-service'


When /^I create an ExoRelay instance .*: "([^"]*)"$/, (code) ->
  eval livescript.compile("@exo-relay = #{code}", bare: yes, header: no)


When /^I send a .*message/, (code) ->
  code = code.replace /\bexo-relay\b/, '@exo-relay'
  @message-id = eval livescript.compile(code, bare: yes, header: no)
  expect(@message-id).to.not.be.undefined


When /^I try to add another handler for that message$/, ->
  try
    @exo-relay.register-handler 'hello', ->
  catch
    @error = e.message


When /^I(?: try to)? register/, (code) ->
  try
    eval livescript.compile("@#{code}", bare: yes, header: no)
  catch
    @error = e.message


When /^I take it online$/, (done) ->
  @exo-relay
    ..connect!
    ..on 'online', ~>
      wait-until (~> @exocom.received-messages.length), 10, ~>
        if @exocom.received-messages |> any (.name is "exocom.register-service")
          @exocom.reset!
        done!


When /^I try to take it online$/, (done) ->
  @exo-relay
    ..connect!
    ..on 'error', (@error) ~>
      done!


When /^receiving the "([^"]*)" message with payload "([^"]*)" as a reply to the "(?:[^"]*)" message$/, (message-name, payload, done) ->
  wait-until (~> @exocom.received-messages.length), 1, ~>
    exocom-data =
      service: 'test-service'
      name: message-name
      payload: payload
      id: '123'
      response-to: @exocom.received-messages[0].id
    @exocom
      ..reset!
      ..send exocom-data
    done!


When /^receiving this message:$/, (request-data) ->
  eval livescript.compile "data = {\n#{request-data}\n}", bare: yes, header: no
  data.message-id = data.id
  data.service = 'test-service'
  @exocom.send data


When /^running this multi\-level request:$/, (code) ->
  done = @done = sinon.stub!
  eval livescript.compile code.replace(/\bexo-relay\b/g, '@exo-relay'), bare: yes, header: no


When /^sending the message:$/, (code) ->
  eval livescript.compile "@message-id = @#{code}", bare: yes, header: no


When /^the reply arrives via this message:$/, (request-data) ->
  rendered = ejs.render request-data, request_uuid: @message-id
  eval livescript.compile "data = {\n#{rendered}\n}", bare: yes, header: no
  data.service = 'test-service'
  @exocom.send data


When /^trying to create an ExoRelay without providing the ExoCom port$/, ->
  try
    @exo-relay = new ExoRelay do
      exocom-host: 'localhost'
      role: 'test-service'
  catch
    @error = e.message


When /^trying to create an ExoRelay without providing the ExoCom host$/, ->
  try
    @exo-relay = new ExoRelay do
      exocom-port: 4100
      role: 'test-service'
  catch
    @error = e.message


When /^trying to send/, (code) ->
  try
    eval livescript.compile "@#{code}", bare: yes, header: no
  catch
    @error = e.message
