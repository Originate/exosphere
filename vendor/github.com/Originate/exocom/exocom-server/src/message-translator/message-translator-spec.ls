require! {
  'chai' : {expect}
  './message-translator' : MessageTranslator
}


describe 'MessageTranslator', ->


  before-each ->
    @message-translator = new MessageTranslator


  describe 'internal-message-name', (...) ->

    it 'translates the given message to the internal format of the given sender', ->
      service =
        internal-namespace: 'text-snippets'
      result = @message-translator.internal-message-name 'tweets.create', for: service
      expect(result).to.eql 'text-snippets.create'


    it 'does not translate the given message if the recipient has the same internal namespace as the message', ->
      service =
        internal-namespace: 'users'
      result = @message-translator.internal-message-name 'users.create', for: service
      expect(result).to.eql 'users.create'



    it 'does not translate the given message if it is not in a translatable format', ->
      result = @message-translator.internal-message-name 'foo bar', for: {}
      expect(result).to.eql 'foo bar'


  describe 'public-message-name', (...) ->

    it "does not convert messages that don't match the format", ->
      result = @message-translator.public-message-name do
        internal-message-name: 'foo bar'
        client-name: 'tweets'
        internal-namespace: 'text-snippets'
      expect(result).to.eql 'foo bar'

    it 'does not convert messages that have the same internal and external namespace', ->
      result = @message-translator.public-message-name do
        internal-message-name: 'users.create'
        client-name: 'users'
        internal-namespace: 'users'
      expect(result).to.eql 'users.create'

    it 'does not convert messages if the service has no internal namespace', ->
      result = @message-translator.public-message-name do
        internal-message-name: 'users.create'
        client-name: 'users'
        internal-namespace: ''
      expect(result).to.eql 'users.create'


    it 'converts messages into the external namespace of the service', ->
      result = @message-translator.public-message-name do
        internal-message-name: 'text-snippets.create'
        client-name: 'tweets'
        internal-namespace: 'text-snippets'
      expect(result).to.eql 'tweets.create'
