const {expect} = require('chai'),
      dimConsole = require('dim-console'),
      ExoComMock = require('exocom-mock'),
      ExoService = require('exoservice'),
      jsDiff = require('jsdiff-console'),
      lowercaseKeys = require('lowercase-keys'),
      N = require('nitroglycerin'),
      portReservation = require('port-reservation'),
      HttpRecorder = require('record-http'),
      request = require('request'),
      {waitUntil} = require('wait')


module.exports = function() {

  this.Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(N( (exocomPort) => {
      this.exocomPort = exocomPort
      this.exocom = new ExoComMock()
      this.exocom.listen(this.exocomPort, done)
    }))
  })


  this.Given(/^an instance of this service$/, function(done) {
    portReservation.getPort(N( (servicePort) => {
      this.servicePort = servicePort
      this.exocom.registerService({ name: '_____serviceName_____',
                                    port: this.servicePort })
      this.process = new ExoService({ serviceName: '_____serviceName_____',
                                      exocomPort: this.exocom.port,
                                      exorelayPort: this.servicePort })
      this.process.listen()
      this.process.on('online', () => done())
    }))
  })


  this.Given(/^the service contains the _____serviceName_____s:$/, function(table, done) {
    _____serviceName_____s = []
    for (record of table.hashes()) {
      _____serviceName_____s.push(lowercaseKeys(record))
    }
    this.exocom.sendMessage({ service: '_____serviceName_____',
                              name: '_____serviceName_____.create-many',
                              payload: _____serviceName_____s })
    this.exocom.waitUntilReceive(done)
  })



  this.When(/^sending the message "([^"]*)"$/, function(message) {
    this.exocom.sendMessage({ service: '_____serviceName_____',
                              name: message })
  })


  this.When(/^sending the message "([^"]*)" with the payload:$/, function(message, payload, done) {
    this.fillIn_____serviceName@camelcase_____Ids(payload, (filledPayload) => {
      this.exocom.sendMessage({ service: '_____serviceName_____',
                                name: message,
                                payload: JSON.parse(filledPayload) })
      done()
    })
  })



  this.Then(/^the service contains no _____serviceName_____s$/, function(done) {
    this.exocom.sendMessage({ service: '_____serviceName_____',
                              name: '_____serviceName_____.list' })
    this.exocom.waitUntilReceive( () => {
      expect(this.exocom.receivedMessages()[0].payload.count).to.equal(0)
      done()
    })
  })


  this.Then(/^the service now contains the _____serviceName_____s:$/, function(table, done) {
    this.exocom.sendMessage({ service: '_____serviceName_____', name: '_____serviceName_____.list' })
    this.exocom.waitUntilReceive( () => {
      actual_____serviceName@camelcase_____s = this.removeIds(this.exocom.receivedMessages()[0].payload._____serviceName_____s)
      expected_____serviceName@camelcase_____s = []
      for (let _____serviceName_____ of table.hashes()) {
        expected_____serviceName@camelcase_____s.push(lowercaseKeys(_____serviceName_____))
      }
      jsDiff(actual_____serviceName@camelcase_____s,
             expected_____serviceName@camelcase_____s,
             done)
    })
  })


  this.Then(/^the service replies with "([^"]*)" and the payload:$/, function(message, payload, done) {
    var expectedPayload = null
    eval(`expectedPayload = ${payload}`)
    this.exocom.waitUntilReceive( () => {
      actualPayload = this.exocom.receivedMessages()[0].payload
      jsDiff(this.removeIds(actualPayload),
             this.removeIds(expectedPayload),
             done)
    })
  })

}
