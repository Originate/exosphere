require! {
  'chai' : {expect}
  '../features/support/text-tools' : {ascii}
}


describe 'text-tools', ->

  describe 'ascii', (...) ->

    it 'returns the ascii value of the given character', ->
      expect(ascii 'a').to.equal 97


    it 'returns the ascii value of the given string', ->
      expect(ascii 'aa').to.equal 194
