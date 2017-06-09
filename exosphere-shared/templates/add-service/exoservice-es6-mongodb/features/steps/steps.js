const {expect} = require('chai'),
      defineSupportCode = require('cucumber').defineSupportCode,
      dimConsole = require('dim-console'),
      ExoComMock = require('exocom-mock'),
      ExoService = require('exoservice'),
      jsDiff = require('jsdiff-console'),
      lowercaseKeys = require('lowercase-keys'),
      N = require('nitroglycerin'),
      portReservation = require('port-reservation'),
      request = require('request'),
      wait = require('wait')


defineSupportCode(function({Given, When, Then}) {

  Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(N( (exocomPort) => {
      this.exocomPort = exocomPort
      this.exocom = new ExoComMock()
      this.exocom.listen(this.exocomPort, done)
    }))
  })


  Given(/^an instance of this service$/, function(done) {
    this.process = new ExoService({ role: '_____serviceRole_____',
                                    exocomHost: 'localhost',
                                    exocomPort: this.exocomPort })
    this.process.connect()
    this.process.on('online', () => wait.wait(10, done))
  })


  Given(/^the service contains the _____modelName_____s:$/, function(table, done) {
    _____modelName_____s = []
    for (record of table.hashes()) {
      _____modelName_____s.push(lowercaseKeys(record))
    }
    this.exocom.send({ service: '_____serviceRole_____',
                       name: '_____modelName_____.create-many',
                       payload: _____modelName_____s })
    this.exocom.onReceive(done)
  })



  When(/^receiving the message "([^"]*)"$/, function(message) {
    this.exocom.send({ service: '_____serviceRole_____',
                       name: message })
  })


  When(/^receiving the message "([^"]*)" with the payload:$/, function(message, payload, done) {
    this.fillIn_____modelNameCamelcase_____Ids(payload, (filledPayload) => {
      this.exocom.send({ service: '_____serviceRole_____',
                         name: message,
                         payload: JSON.parse(filledPayload) })
      done()
    })
  })



  Then(/^the service contains no _____modelName_____s$/, function(done) {
    this.exocom.send({ service: '_____serviceRole_____',
                       name: '_____modelName_____.list' })
    this.exocom.onReceive( () => {
      expect(this.exocom.receivedMessages[0].payload.count).to.equal(0)
      done()
    })
  })


  Then(/^the service now contains the _____modelName_____s:$/, function(table, done) {
    this.exocom.send({ service: '_____serviceRole_____', name: '_____modelName_____.list' })
    this.exocom.onReceive( () => {
      actual_____modelNameCamelcase_____s = this.removeIds(this.exocom.receivedMessages[0].payload)
      expected_____modelNameCamelcase_____s = []
      for (let _____modelName_____ of table.hashes()) {
        expected_____modelNameCamelcase_____s.push(lowercaseKeys(_____modelName_____))
      }
      jsDiff(actual_____modelNameCamelcase_____s,
             expected_____modelNameCamelcase_____s,
             done)
    })
  })


  Then(/^the service replies with "([^"]*)" and the payload:$/, function(message, payload, done) {
    var expectedPayload = null
    eval(`expectedPayload = ${payload}`)
    this.exocom.onReceive( () => {
      expect(this.exocom.receivedMessages[0].name).to.equal(message)
      actualPayload = this.exocom.receivedMessages[0].payload
      jsDiff(this.removeIds(actualPayload),
             this.removeIds(expectedPayload),
             done)
    })
  })

});
