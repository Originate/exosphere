require! {
  'cucumber': {defineSupportCode}
  '../../../exocom-mock' : ExoComMock
  'http'
  'nitroglycerin': N
  'port-reservation'
  'wait' : {wait}
}


defineSupportCode ({Given}) ->

  Given /^an ExoCom instance$/, (done) ->
    port-reservation.get-port N (@exocom-port) ~>
      @exocom = new ExoComMock
        ..listen @exocom-port
      done!

  Given /^an ExoCom instance running at port (\d+)$/, (@exocom-port) ->
    @exocom = new ExoComMock
      ..listen @exocom-port


  Given /^an instance of the "([^"]*)" service$/, (@role, done) ->
    @create-exoservice-instance {@role, @exocom-port}, ~>
      @remove-register-service-message @exocom, done


  Given /^ports (\d+) and (\d+) are used$/, (port1, port2, done) ->
    # Note: this is due to a Cucumber-JS issue where cleanup methods aren't async.
    # So we have to let all remaining messages in the event queue be processed here
    # so that any code that releases ports has actually been executed.
    wait 100, ~>
      handler = (_, res) -> res.end 'existing server'
      @server1 = http.create-server(handler).listen 3000, 'localhost', ~>
        @server2 = http.create-server(handler).listen 3001, 'localhost', done
