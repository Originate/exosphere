require! {
  'async'
  'chai'
  'jsdiff-console'
  'lowercase-keys'
  'prelude-ls' : {filter, map}
  'request'
  'sinon-chai'
  'wait' : {wait-until}
}
expect = chai.expect
chai.use sinon-chai


module.exports = ->

  @Then /^ExoCom now knows about these services:$/ (table) ->
    async.each table.raw!,
               (~> @exocom.wait-until-knows-service &1, &2)


  @Then /^ExoComMock makes the request:$/, (table, done) ->
    @verify-exocom-received-request table.rows-hash!, done


  @Then /^I can close it without errors$/, (done) ->
    @exocom.close done


  @Then /^I get the error "([^"]*)"$/, (expected-error) ->
    expect(@error.message).to.equal expected-error


  @Then /^it calls the given callback$/, (done) ->
    wait-until (~> @call-received.called), done


  @Then /^it calls the given callback right away$/, (done) ->
    wait-until (~> @exocom.received-messages.length), 1, ~>
      expect(@call-received).to.have.been.called
      done!


  @Then /^it doesn't call the given callback right away$/, ->
    expect(@call-received).to.not.have.been.called


  @Then /^it has received no messages/, ->
    expect(@exocom.received-messages).to.be.empty


  @Then /^it has received the messages/, (table, done) ->
    wait-until (~> @exocom.received-messages.length > 1), 10, ~>
      expected-messages = table.hashes! |> map lowercase-keys
      service-messages = filter (.name is not "exocom.register-service"), @exocom.received-messages
      jsdiff-console service-messages, expected-messages, done


  @Then /^it is no longer listening$/, (done) ->
    request-data =
      url: "http://localhost:#{@exocom-port}/send/foo"
      method: 'POST'
      body:
        payload: ''
      json: yes
    request request-data, (err) ~>
      expect(err).to.not.be.undefined
      expect(err.message).to.equal "connect ECONNREFUSED 127.0.0.1:#{@exocom-port}"
      done!
