require! {
  './websocket-connector' : WebSocketConnector
  'chai' : {expect}
  'sinon'
}

describe 'WebSocketConnector', ->

  before-each ->
    @websocket-connector = new WebSocketConnector exocom-host: 'localhost', exocom-port: 4100, role: 'test'
      ..connect!
      ..on 'error', (@error) ~>

  after-each (done) ->
    @websocket-connector.close done


  describe 'reply-method-for', (...) ->

    before-each ->
      @websocket-connector.send = sinon.stub!
      @reply-method = @websocket-connector.reply-method-for '123', '1'

    it 'returns a function that calls @send prebound to the response id', ->
      expect(@reply-method).to.be.a 'function'
      @reply-method 'reply-message', 'payload'
      expect(@websocket-connector.send.first-call.args).to.eql [ 'reply-message', 'payload', {response-to: '123', "sessionId": '1'} ]


    context 'missing id', (...) ->

      before-each ->
        @websocket-connector.reply-method-for null

      it 'emits an error', (done) ->
        expect(@error.message).to.eql 'WebSocketConnector.replyMethodFor needs an id'
        done!
