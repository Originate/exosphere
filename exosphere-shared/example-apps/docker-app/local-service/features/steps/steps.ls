require! {
  'chai' : {expect}
  'dim-console'
  'exocom-mock' : ExoComMock
  'exoservice' : ExoService
  'jsdiff-console'
  'livescript'
  'lowercase-keys'
  'nitroglycerin' : N
  'port-reservation'
  'request'
  'wait' : {wait-until, wait}
}


module.exports = ->

  @Given /^an ExoCom server$/, (done) ->
    port-reservation.get-port N (@exocom-port) ~>
      @exocom = new ExoComMock
        ..listen @exocom-port, done


  @Given /^an instance of this service$/, (done) ->
    @process = new ExoService service-name: 'local-service', exocom-port: @exocom-port, exocom-host: 'localhost'
      ..connect!
      ..on 'online', -> wait 10, done


  @Given /^the service contains the tests:$/, (table, done) ->
    tests = [lowercase-keys(record) for record in table.hashes!]
    @exocom
      ..send service: 'local-service', name: 'test.create-many', payload: tests
      ..on-receive done



  @When /^sending the message "([^"]*)"$/, (message) ->
    @exocom.send service: 'local-service', name: message


  @When /^sending the message "([^"]*)" with the payload:$/, (message, payload, done) ->
    @fill-in-test-ids payload, (filled-payload) ~>
      if filled-payload[0] is '['   # payload is an array
        eval livescript.compile "payload-json = #{filled-payload}", bare: true, header: no
      else                          # payload is a hash
        eval livescript.compile "payload-json = {\n#{filled-payload}\n}", bare: true, header: no
      @exocom.send service: 'local-service', name: message, payload: payload-json
      done!



  @Then /^the service contains no tests$/, (done) ->
    @exocom
      ..send service: 'local-service', name: 'test.list'
      ..on-receive ~>
        expect(@exocom.received-messages[0].payload.count).to.equal 0
        done!


  @Then /^the service now contains the tests:$/, (table, done) ->
    @exocom
      ..send service: 'local-service', name: 'test.list'
      ..on-receive ~>
        actual-tests = @remove-ids @exocom.received-messages[0].payload
        expected-tests = [lowercase-keys(test) for test in table.hashes!]
        jsdiff-console actual-tests, expected-tests, done


  @Then /^the service replies with "([^"]*)" and the payload:$/, (message, payload, done) ->
    expected-payload = eval livescript.compile payload, bare: true
    @exocom.on-receive ~>
      expect(@exocom.received-messages[0].name).to.equal message
      actual-payload = @exocom.received-messages[0].payload
      jsdiff-console @remove-ids(actual-payload), @remove-ids(expected-payload), done
