const {expect} = require('chai'),
      dimConsole = require('dim-console'),
      ExoComMock = require('exocom-mock'),
      ExoService = require('exoservice'),
      jsDiff = require('jsdiff-console'),
      lowercaseKeys = require('lowercase-keys'),
      N = require('nitroglycerin'),
      portReservation = require('port-reservation'),
      request = require('request'),
      {waitUntil} = require('wait')


module.exports = function() {

  this.Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(N( (exocomPort) => {
      this.exocomPort = exocomPort
      this.exocom = new ExoComMock()
      this.exocom.listen(this.exocomPort)
      done()
    }))
  })


  this.Given(/^an instance of this service$/, function(done) {
    portReservation.getPort(N( (servicePort) => {
      this.servicePort = servicePort
      this.exocom.registerService({ name: '_____serviceName_____',
                                    port: this.servicePort })
      this.process = new ExoService({ serviceName: '_____serviceName_____',
                                      exocomPort: this.exocom.pullSocketPort,
                                      exorelayPort: this.servicePort })
      this.process.listen()
      this.process.on('online', () => done())
    }))
  })


  this.Given(/^the service contains the _____modelName_____s:$/, function(table, done) {
    _____modelName_____s = []
    for (record of table.hashes()) {
      _____modelName_____s.push(lowercaseKeys(record))
    }
    this.exocom.send({ service: '_____serviceName_____',
                              name: '_____modelName_____.create-many',
                              payload: _____modelName_____s })
    this.exocom.onReceive(done)
  })



  this.When(/^receiving the message "([^"]*)"$/, function(message) {
    this.exocom.send({ service: '_____serviceName_____',
                              name: message })
  })


  this.When(/^receiving the message "([^"]*)" with the payload:$/, function(message, payload, done) {
    this.fillIn_____modelName@camelcase_____Ids(payload, (filledPayload) => {
      this.exocom.send({ service: '_____serviceName_____',
                                name: message,
                                payload: JSON.parse(filledPayload) })
      done()
    })
  })



  this.Then(/^the service contains no _____modelName_____s$/, function(done) {
    this.exocom.send({ service: '_____serviceName_____',
                              name: '_____modelName_____.list' })
    this.exocom.onReceive( () => {
      expect(this.exocom.receivedMessages[0].payload.count).to.equal(0)
      done()
    })
  })


  this.Then(/^the service now contains the _____modelName_____s:$/, function(table, done) {
    this.exocom.send({ service: '_____serviceName_____', name: '_____modelName_____.list' })
    this.exocom.onReceive( () => {
      actual_____modelName@camelcase_____s = this.removeIds(this.exocom.receivedMessages[0].payload)
      expected_____modelName@camelcase_____s = []
      for (let _____modelName_____ of table.hashes()) {
        expected_____modelName@camelcase_____s.push(lowercaseKeys(_____modelName_____))
      }
      jsDiff(actual_____modelName@camelcase_____s,
             expected_____modelName@camelcase_____s,
             done)
    })
  })


  this.Then(/^the service replies with "([^"]*)" and the payload:$/, function(message, payload, done) {
    var expectedPayload = null
    eval(`expectedPayload = ${payload}`)
    this.exocom.onReceive( () => {
      actualPayload = this.exocom.receivedMessages[0].payload
      jsDiff(this.removeIds(actualPayload),
             this.removeIds(expectedPayload),
             done)
    })
  })

}
