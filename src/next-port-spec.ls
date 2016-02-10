require! {
  'chai' : {expect}
  './next-port'
}


describe 'next-port', (...) ->

  it 'returns consecutive ports', (done) ->
    next-port (port1) ->
      next-port (port2) ->
        expect(port1).to.not.equal port2
        done!
