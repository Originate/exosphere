require! {
  'cucumber': {defineSupportCode}
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


defineSupportCode ({Given, When, Then}) ->

  Given /^an ExoCom server$/, (done) ->
    port-reservation.get-port N (@exocom-port) ~>
      @exocom = new ExoComMock
        ..listen @exocom-port, done


  Given /^an instance of this service$/, (done) ->
    @process = new ExoService role: '_____serviceRole_____', exocom-port: @exocom-port, exocom-host: 'localhost'
      ..connect!
      ..on 'online', -> wait 10, done


  Given /^the service contains the _____modelName_____s:$/, (table, done) ->
    _____modelName_____s = [lowercase-keys(record) for record in table.hashes!]
    @exocom
      ..send service: '_____serviceRole_____', name: '_____modelName_____.create-many', payload: _____modelName_____s
      ..on-receive done


  When /^sending the message "([^"]*)"$/, (message) ->
    @exocom.send service: '_____serviceRole_____', name: message


  When /^sending the message "([^"]*)" with the payload:$/, (message, payload, done) ->
    @fill-in-_____modelName_____-ids payload, (filled-payload) ~>
      if filled-payload[0] is '['   # payload is an array
        eval livescript.compile "payload-json = #{filled-payload}", bare: true, header: no
      else                          # payload is a hash
        eval livescript.compile "payload-json = {\n#{filled-payload}\n}", bare: true, header: no
      @exocom.send service: '_____serviceRole_____', name: message, payload: payload-json
      done!



  Then /^the service contains no _____modelName_____s$/, (done) ->
    @exocom
      ..send service: '_____serviceRole_____', name: '_____modelName_____.list'
      ..on-receive ~>
        expect(@exocom.received-messages[0].payload.count).to.equal 0
        done!


  Then /^the service now contains the _____modelName_____s:$/, (table, done) ->
    @exocom
      ..send service: '_____serviceRole_____', name: '_____modelName_____.list'
      ..on-receive ~>
        actual_____modelName_____s = @remove-ids @exocom.received-messages[0].payload
        expected_____modelName_____s = [lowercase-keys(_____modelName_____) for _____modelName_____ in table.hashes!]
        jsdiff-console actual_____modelName_____s, expected_____modelName_____s, done


  Then /^the service replies with "([^"]*)" and the payload:$/, (message, payload, done) ->
    expected-payload = eval livescript.compile payload, bare: true
    @exocom.on-receive ~>
      expect(@exocom.received-messages[0].name).to.equal message
      actual-payload = @exocom.received-messages[0].payload
      jsdiff-console @remove-ids(actual-payload), @remove-ids(expected-payload), done
