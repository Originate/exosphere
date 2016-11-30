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
      this.exocom.registerService({ name: 'todo',
                                    port: this.servicePort })
      this.process = new ExoService({ serviceName: 'todo',
                                      exocomPort: this.exocom.port,
                                      exorelayPort: this.servicePort })
      this.process.listen()
      this.process.on('online', () => done())
    }))
  })


  this.Given(/^the service contains the todos:$/, function(table, done) {
    todos = []
    for (record of table.hashes()) {
      todos.push(lowercaseKeys(record))
    }
    this.exocom.sendMessage({ service: 'todo',
                              name: 'todo.create-many',
                              payload: todos })
    this.exocom.waitUntilReceive(done)
  })



  this.When(/^sending the message "([^"]*)"$/, function(message) {
    this.exocom.sendMessage({ service: 'todo',
                              name: message })
  })


  this.When(/^sending the message "([^"]*)" with the payload:$/, function(message, payload, done) {
    this.fillIntodoIds(payload, (filledPayload) => {
      this.exocom.sendMessage({ service: 'todo',
                                name: message,
                                payload: JSON.parse(filledPayload) })
      done()
    })
  })



  this.Then(/^the service contains no todos$/, function(done) {
    this.exocom.sendMessage({ service: 'todo',
                              name: 'todo.list' })
    this.exocom.waitUntilReceive( () => {
      expect(this.exocom.receivedMessages()[0].payload.count).to.equal(0)
      done()
    })
  })


  this.Then(/^the service now contains the todos:$/, function(table, done) {
    this.exocom.sendMessage({ service: 'todo', name: 'todo.list' })
    this.exocom.waitUntilReceive( () => {
      actualtodos = this.removeIds(this.exocom.receivedMessages()[0].payload)
      expectedtodos = []
      for (let todo of table.hashes()) {
        expectedtodos.push(lowercaseKeys(todo))
      }
      jsDiff(actualtodos,
             expectedtodos,
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
