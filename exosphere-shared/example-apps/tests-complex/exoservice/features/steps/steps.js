defineSupportCode = require('cucumber').defineSupportCode;
ExoComMock = require('exocom-mock')
ExoService = require('exoservice')
expect = require('chai').expect
fs = require('fs')
N = require('nitroglycerin')
portReservation = require('port-reservation')
yaml = require('js-yaml')
wait = require('wait')


const serviceConfig = yaml.safeLoad(fs.readFileSync('service.yml'), 'utf8')


defineSupportCode(function({Given, When, Then}) {

  Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(N( exocomPort => {
      this.exocomPort = exocomPort
      this.exocom = new ExoComMock()
      this.exocom.listen(exocomPort, done)
    }))
  })


  Given(/^an instance of this service$/, function(done) {
    this.process = new ExoService({  role: serviceConfig.type,
                                     exocomPort: this.exocomPort,
                                     exocomHost: 'localhost' })
    this.process.connect()
    this.process.on('online', () => wait.wait(10, done))
  })


  When(/^receiving the "([^"]*)" command$/, function(commandName) {
    this.exocom.reset()
    this.exocom.send({ service: serviceConfig.type,
                              name: commandName })
  })


  Then(/^this service replies with a "([^"]*)" message/, function(expectedMessageName, done) {
    this.exocom.onReceive( () => {
      const receivedMessages = this.exocom.receivedMessages
      expect(receivedMessages).to.have.length(1)
      expect(receivedMessages[0].name).to.equal(expectedMessageName)
      done()
    })
  })

});
