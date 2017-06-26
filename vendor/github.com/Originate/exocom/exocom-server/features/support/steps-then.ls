require! {
  'livescript'
  'wait' : {wait}
}


module.exports = ->

  @Then /^ExoCom broadcasts the message "([^"]*)" to the "([^"]+)" service$/, (message-name, client-name, done) ->
    @verify-sent-calls {client-name, message-name, id: @last-sent-message-id}, done


  @Then /^ExoCom broadcasts the reply "([^"]*)" to the "([^"]+)" service$/, (message-name, client-name, done) ->
    @verify-sent-calls {client-name, message-name, response-to: '111'}, done


  @Then /^ExoCom now knows about these service instances:$/ (table, done) ->
    services = {}
    for row in table.hashes!
      services[row['CLIENT NAME']] =
        client-name: row['CLIENT NAME']
        service-type: row['SERVICE TYPE']
        internal-namespace: row['INTERNAL NAMESPACE']
    @verify-service-setup services, done


  @Then /^ExoCom signals "([^"]*)"$/, (message-name, done) ->
    @verify-exocom-signaled-string message-name, done


  @Then /^ExoCom signals that this message was sent$/, (done) ->
    @verify-exocom-broadcasted-message message: @last-sent-message-name, done


  @Then /^ExoCom signals that this reply is sent from the ([^ ]+) to the (.+)$/, (sender, receiver, done) ->
    @verify-exocom-broadcasted-message message: @last-sent-message-name, sender: sender, receivers: [receiver], response-to: '111', done


  @Then /^ExoCom signals that this reply was sent$/, (done) ->
    @verify-exocom-broadcasted-reply @last-sent-message-name, done


  @Then /^ExoCom signals the error "([^"]*)"$/, (error-text, done) ->
    @process.wait error-text, done


  @Then /^it aborts with the message "([^"]*)"$/, (text, done) ->
    @verify-abort-with-notification text, done


  @Then /^it has this routing table:$/, (table, done) ->
    expected-routes = {}
    for row in table.hashes!
      eval livescript.compile "receiver-json = {#{row.RECEIVERS}}", bare: yes, header: no
      expected-routes[row.MESSAGE] =
        receivers: [receiver-json]
    @verify-routing-setup expected-routes, done


  @Then /^it opens a port at (\d+)$/, (+port, done) ->
    @verify-listening-at-ports port, done


  @Then /^it opens an HTTP listener at port (\d+)$/, (+port, done) ->
    @verify-listening-at-ports port, done
