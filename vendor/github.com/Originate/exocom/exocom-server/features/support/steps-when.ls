require! {
  'livescript'
  '../support/mock-service' : MockService
  'prelude-ls' : {pairs-to-obj}
  'wait' : {wait}
}


module.exports = ->

  @When /^a new "([^"]*)" service instance registers itself with it via the message:$/ (client-name, table, done) ->
    table-data = table.raw! |> pairs-to-obj
    payload = table-data.PAYLOAD |> JSON.parse
    (@service-mocks or= {})[client-name] = new MockService name: table-data.NAME, port: @exocom-port
    @service-mocks[client-name].connect {payload}, ->
      wait 200, done


  @When /^(I try )?starting ExoCom at port (\d+)$/, (!!expect-error, +port, done) ->
    @run-exocom-at-port port, expect-error, done


  @When /^I( try to)? run ExoCom$/, (!!expect-error, done) ->
    @run-exocom-at-port null, expect-error, done


  @When /^requesting the routing information$/, (done) ->
    @get-routing-information (@routing-information) ~>
      done!


  @When /^sending the service configuration:$/, (config-text, done) ->
    eval livescript.compile "config = \n#{config-text.replace /^/gm, '  '}", bare: yes, header: no
    @set-service-landscape config, done


  @When /^the "([^"]*)" service goes offline$/ (client-name, done) ->
    @service-mocks[client-name].close!
    wait 200, done


  @When /^the "([^"]+)" service sends "([^"]*)"$/, (service-name, message-name, done) ->
    @last-sent-message = message-name
    @service-sends-message {service-name, message-name}, done


  @When /^the "([^"]+)" service sends "([^"]*)" in reply to "([^"]*)"$/, (service-name, reply-message, response-to, done) ->
    @last-sent-message-name = reply-message
    @service-sends-reply {service-name, message-name: reply-message, response-to}, done
