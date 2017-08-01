require! {
  '../..' : MockExoCom
  'chai'
  'cucumber': {defineSupportCode}
  'nitroglycerin': N
  'port-reservation'
  'wait' : {wait-until}
}


defineSupportCode ({Given}) ->

  Given /^a known "([^"]*)" service$/, (name, done) ->
    @create-named-websocket-endpoint {name, @exocom-port}, done


  Given /^an ExoComMock instance$/, (done) ->
    @exocom = new MockExoCom
    port-reservation.get-port N (@exocom-port) ~>
      @exocom.listen @exocom-port, done


  Given /^somebody sends it a message$/, (done) ->
    @create-websocket-endpoint {@exocom-port}, ~>
      old-length = @exocom.received-messages.length
      @service-send-message name: \foo, payload: '', id: \123
      wait-until (~> @exocom.received-messages.length > old-length), 1, done


  Given /^somebody sends it a "([^"]*)" message with payload "([^"]*)"$/, (name, payload, done) ->
    @create-websocket-endpoint {@exocom-port}, ~>
      old-length = @exocom.received-messages.length
      @service-send-message name: name, payload: payload, id: \123
      wait-until (~> @exocom.received-messages.length > old-length), 1, done
