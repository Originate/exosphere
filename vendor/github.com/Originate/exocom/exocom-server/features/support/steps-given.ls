require! {
  'http'
  'nitroglycerin' : N
  'port-reservation'
  'wait' : {wait}
}


module.exports = ->

  @Given /^a new "([^"]*)" service$/ (client-name, done) ->
    @create-mock-service-at-port {client-name, port: @exocom-port}, done


  @Given /^a running "([^"]*)" instance$/, (client-name, done) ->
    @create-mock-service-at-port {client-name, port: @exocom-port}, ->
      wait 200, done


  @Given /^an ExoCom instance$/, (done) ->
    port-reservation
      ..base-port = 5000
      ..get-port N (@exocom-port) ~>
        @create-exocom-instance port: @exocom-port, done


  @Given /^an ExoCom instance configured with the routes:?$/, (service-routes, done) ->
    service-routes = service-routes |> (.replace /\s/g, '') |> (.replace /"/g, '')
    port-reservation
      ..base-port = 5000
      ..get-port N (@exocom-port) ~>
        @create-exocom-instance port: @exocom-port, service-routes: service-routes, done


  @Given /^an ExoCom instance managing the service landscape:$/ (table, done) ->
    port-reservation
      ..base-port = 5000
      ..get-port N (@exocom-port) ~>
        @create-exocom-instance port: @exocom-port, ~>
          @create-mock-service-at-port {client-name: service.NAME, port: @exocom-port}, ->
          done!


  @Given /^another service already uses port (\d+)$/, (+port, done) ->
    wait 100, ~>   # to let the previous server shut down
      handler = (_, res) -> res.end 'existing server'
      @existing-server = http.create-server(handler).listen port, done


  @Given /^the "([^"]+)" service sends "([^"]*)" with id "([^"]*)"$/, (service-name, message-name, id, done) ->
    @last-sent-message-name = message-name
    @last-sent-message-id = id
    @service-sends-message {service-name, message-name, id}, done
