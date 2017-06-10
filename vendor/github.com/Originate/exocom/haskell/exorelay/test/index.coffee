ExoComMock = require 'exocom-mock'

exocom = new ExoComMock


exocom.listen 4100
exocom.registerService {name: 'exorelay-hs', port: 4001 }

exocom.onReceive =>
  received = exocom.receivedMessages[0]
  if received.name is 'hello'
    received.service = 'exorelay-hs'
    exocom.send received
  else if received.name is 'reply'
    received.service = 'exorelay-hs'
    received.responseTo = received.id
    exocom.send received
    setTimeout (->
      exocom.send {sender: 'exorelay-hs', service: 'exorelay-hs', payload: 'payload', name: 'needReply', messageId:'1234'}), 1000
  else if received.responseTo is '1234'
    console.log 'received reply'
    received.service = 'exorelay-hs'
    received.responseTo = received.id
    exocom.send received
  else if received.name is 'listenReply'
    received.service = 'exorelay-hs'
    exocom.send received
  else
    console.log received
