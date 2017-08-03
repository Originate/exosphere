require! {
  'livescript'
  'sinon'
  'wait' : {wait-until}
  '../..' : ExoRelay
  'exocom-mock': MockExoCom
  'prelude-ls' : {any}
  'cucumber': {Given}
}


Given /^a "([^"]*)" message$/, (message-name) ->
  global[message-name] = sinon.stub!


Given /^a new ExoRelay instance connecting to port (\d+)$/ (@exocom-port) ->
  @exo-relay = new ExoRelay {exocom-host: 'localhost', @exocom-port, role: @role = 'test-service'}


Given /^an ExoRelay instance$/, (done) ->
  @exo-relay = new ExoRelay {exocom-host: 'localhost', @exocom-port, role: @role = \test-service}
    ..connect!
    ..on 'online', ~>
      wait-until (~> @exocom.received-messages.length), 10, ~>
        if @exocom.received-messages |> any (.name is "exocom.register-service")
          @exocom.reset!
        done!
    ..on 'error', (@error) ~>


Given /^an ExoRelay instance called "([^"]*)" running inside the "([^"]*)" service$/, (instance-name, @role, done) ->
  @exo-relay = new ExoRelay {exocom-host: 'localhost', @role, @exocom-port}
    ..connect!
    ..on 'online', ~>
      wait-until (~> @exocom.received-messages.length), 10, ~>
        if @exocom.received-messages |> any (.name is "exocom.register-service")
          @exocom.reset!
        done!
    ..on 'error', (@error) ~>


Given /^ExoCom runs at port (\d+)$/, (@exocom-port) ->
  @exocom = new MockExoCom
    ..listen @exocom-port


Given /^I register this handler for the "([^"]*)" message:$/, (message-name, code) ->
  code = "handler = #{code}"
  eval livescript.compile code, bare: yes, header: no
  @exo-relay.register-handler message-name, handler


Given /^I try to set up this handler:$/, (code) ->
  try
    eval livescript.compile "@#{code}", bare: yes, header: no
  catch
    @error = e.message


Given /^my ExoRelay instance already has a handler for the message "([^"]*)"$/, (message-name) ->
  @exo-relay.register-handler message-name, ->


Given /^the "([^"]*)" message has this handler:$/, (message-name, handler-code) ->
  eval livescript.compile "@#{handler-code}", bare: yes, header: no
