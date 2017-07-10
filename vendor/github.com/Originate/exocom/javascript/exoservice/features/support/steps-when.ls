require! {
  'livescript'
}


module.exports = ->

  @When /^receiving the( unknown)? "([^"]*)" message$/, (expect-error, message-name) ->
    @exocom
      ..reset!
      ..send {service: @role, name: message-name, expect-error}


  @When /^receiving the( unknown)? "([^"]*)" message with the payload:$/, (expect-error, message-name, payload) ->
    eval livescript.compile "json-payload = {\n#{payload}\n}", bare: yes, header: no
    @exocom
      ..reset!
      ..send {service: @role, name: message-name, payload: json-payload, expect-error}


  @When /^starting a service$/, (done) ->
    @role = 'test service'
    @create-exoservice-instance {@role, @exocom-port}, ~>
      @remove-register-service-message @exocom, done


  @When /^starting a service configured for ExoCom port (\d+)$/, (port, done) ->
    @role = 'test service'
    @create-exoservice-instance {@role, exocom-port: port}, ~>
      @remove-register-service-message @exocom, done


  @When /^starting the "([^"]*)" service$/, (@role, done) ->
    @create-exoservice-instance {@role, @exocom-port}, ~>
      @remove-register-service-message @exocom, done


  @When /^trying to start a service configured for ExoCom port (\d+)$/, (port) ->
    @role = 'test service'
    @create-exoservice-instance {@role, exocom-port: port}
