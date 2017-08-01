require! {
  \chai : {expect}
  \wait : {wait, wait-until}
  \../support/mock-exorelay : MockExorelay
}


World = !->

  @create-websocket-endpoint = (port, done) ->
    | @service  =>  return done!
    @service = new MockExorelay
      ..connect port, done


  @create-named-websocket-endpoint = ({name, exocom-port, registration-message, registration-payload}, done) ->
    @service = new MockExorelay name
      ..connect {exocom-port, registration-message, registration-payload}, ~>
          @exocom.wait-until-knows-service name, done


  @exocom-send-message = ({exocom, service, message-data}) ->
    exocom.send service: service, name: message-data.name, payload: message-data.payload



  @service-send-message = (message-data) ->
    @service.send message-data


  @verify-exocom-received-request = (expected-request, done) ->
    wait-until (~> @service.received-messages.length), 1, ~>
      actual-request = @service.received-messages[0]
      expect(actual-request.name).to.equal expected-request.NAME
      expect(actual-request.payload).to.equal expected-request.PAYLOAD
      done!


module.exports = World
