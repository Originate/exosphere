require! {
  'chai' : {expect}
  'wait' : {wait-until}
}


module.exports = ->

  @Then /^after a while it sends the "([^"]*)" message$/, (reply-message-name, done) ->
    @exocom.on-receive ~>
      received-messages = @exocom.received-messages
      expect(received-messages).to.have.length 1
      expect(received-messages[0].name).to.equal reply-message-name
      done!


  @Then /^after a while it sends the "([^"]*)" message with the textual payload:$/, (reply-message-name, payload-text, done) ->
    @exocom.on-receive ~>
      received-messages = @exocom.received-messages
      expect(received-messages).to.have.length 1
      expect(received-messages[0].name).to.equal reply-message-name
      expect(received-messages[0].payload).to.equal payload-text
      done!


  @Then /^it acknowledges the received message$/, (done) ->
    wait-until (~> @exocom.received-messages.length), done


  @Then /^it can run the "([^"]*)" service$/, (@role, done) ->
    @create-exoservice-instance {@role, @exocom-port}, done


  @Then /^it connects to the ExoCom instance$/, (done) ->
    @exocom.send service: @role, name: '__status' , id: '123'
    wait-until (~> @exocom.received-messages[0]), 1, ~>
      if @exocom.received-messages[0].name is "__status-ok"
        done!


  @Then /^it runs the "([^"]*)" hook$/, (hook-name, done) ->
    @exocom
      ..reset!
      ..send name: 'which-hooks-ran', service: @role
      ..on-receive ~>
        expect(@exocom.received-messages[0].payload).to.eql ['before-all']
        done!
