require! {
  'chai': {expect}
  './docker-helper': DockerHelper
}

describe 'DockerHelper', ->

  describe 'cat-file', ->

    specify 'should parse non-service log message correctly' ->
      expect(0).to.eql(0)
