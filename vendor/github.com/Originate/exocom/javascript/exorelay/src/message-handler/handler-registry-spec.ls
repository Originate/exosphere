require! {
  './handler-registry' : HandlerRegistry
  'chai' : {expect}
  'sinon'
}

describe 'HandlerRegistry', ->

  describe 'handle-command' (_) ->

    it 'leaves the handler after it was processed', ->
      handler-registry = new HandlerRegistry 'test'
      handler = sinon.stub!
      handler-registry.register-handler '1234', handler
      handler-registry.handle-command message-name: '1234'
      expect(handler-registry.has-handler '1234').to.be.true


  describe 'handle-reply' (_) ->

    it 'leaves the handler after it was processed', ->
      handler-registry = new HandlerRegistry 'test'
      handler = sinon.stub!
      handler-registry.register-handler '1234', handler
      handler-registry.handle-reply message-name: '1234'
      expect(handler-registry.has-handler '1234').to.be.true

