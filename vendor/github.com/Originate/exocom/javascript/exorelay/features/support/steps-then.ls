require! {
  'chai' : {expect}
  'ip'
  'jsdiff-console'
  'ejs'
  'livescript'
  'wait' : {wait, wait-until}
  'prelude-ls' : {any}
  'lodash.isequal' : is-equal
}


module.exports = ->

  @Then /^ExoRelay connects to ExoCom$/ (done) ->
    @exo-relay
      ..connect!
    wait-until (~> @exocom.service-sockets[@role]), 1, ~>
      done!


  @Then /^ExoRelay emits an "error" event with the error "([^"]*)"$/, (error-message, done) ->
    wait-until (~> @error), 1, ~>
      expect(@error.message).to.equal error-message
      @error = null
      done!


  @Then /^ExoRelay makes the WebSocket request:$/, (request-data, done) ->
    # Wait until we get some call data, then wait another 50ms to let all the request data fill in
    wait-until (~> @exocom.received-messages.length), 10, ~>
      wait 50, ~>
        rendered = ejs.render request-data, request_uuid: @message-id
        template = livescript.compile "compiled = {\n#{rendered}\n}", bare: yes, header: no
        eval template
        jsdiff-console @exocom.received-messages[0], compiled, done


  @Then /^(?:ExoRelay|it) runs the registered handler, in this example calling "([^"]*)" with "([^"]*)"$/, (message-name, message-argument, done) ->
    wait-until (~> global[message-name].called), 10, done


  @Then /^ExoRelay sends the "([^"]*)" message with payload "([^"]*)"$/, (message-name, payload, done) ->
    wait-until (~> @exocom.received-messages.length), 10, ~>
      wait 50, ~>
        expect(@exocom.received-messages[0].name).to.equal "#{message-name}"
        expect(@exocom.received-messages[0].payload).to.equal payload
        done!


  @Then /^it connects to the given ExoCom host and port$/, (done) ->
    @exocom
      ..send service: @role, name: '__status'
    current-length = @exocom.received-messages.length
    wait-until (~> @exocom.received-messages.length > current-length), 1, ~>
      if @exocom.received-messages |> any (.name is "__status-ok")
        done!


  @Then /^it signals it is online$/, (done) ->
    wait-until (~> @status-code), 1, ~>
      expect(@status-code).to.equal '__status-ok'
      done!


  @Then /^it throws the error "([^"]*)"$/, (expected-error) ->
    expect(@error).to.equal expected-error


  @Then /^my handler calls the "done" method$/, (done) ->
    wait-until (~> @done.called), 10, done


  @Then /^my message handler (?:replies with|sends out) the message:$/ (request-data, done) ->
    # Wait until we get some call data, then wait another 50ms to let all the request data fill in
    wait-until (~> @exocom.received-messages.length), 10, ~>
      wait 50, ~>
        rendered = ejs.render request-data, request_uuid: @exo-relay.websocket-connector.last-sent-id, ip_address: ip.address!
        template = "compiled = {\n#{rendered}\n}"
        compiled-template = livescript.compile template, bare: yes, header: no
        parsed = eval compiled-template
        jsdiff-console @exocom.received-messages[0], parsed, done


  @Then /^the instance has a handler for the message "([^"]*)"$/, (handler1) ->
    expect(@exo-relay.has-handler handler1).to.be.true


  @Then /^the instance has handlers for the messages "([^"]*)" and "([^"]*)"$/, (handler1, handler2) ->
    expect(@exo-relay.has-handler handler1).to.be.true
    expect(@exo-relay.has-handler handler2).to.be.true


  @Then /^the reply handler runs and in this example calls my "([^"]*)" method with "([^"]*)"$/, (message-name, message-args, done) ->
    condition = ~>
      global[message-name].called and is-equal global[message-name].first-call.args, [message-args]
    wait-until condition, 10, done


  @Then /^this instance uses the ExoCom host "([^"]*)" and port (\d+)$/ (host, +port) ->
    expect(@exo-relay.websocket-connector.exocom-port).to.equal port
    expect(@exo-relay.websocket-connector.exocom-host).to.equal host
