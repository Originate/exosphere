ExoComMock = require('exocom-mock')
ExoService = require('exoservice')
expect = require('chai').expect
fs = require('fs')
N = require('nitroglycerin')
portReservation = require('port-reservation')
yaml = require('js-yaml')


const serviceConfig = yaml.safeLoad(fs.readFileSync('service.yml'), 'utf8')


module.exports = function() {

  this.Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(N( exocomPort => {
      this.exocomPort = exocomPort
      this.exocom = new ExoComMock()
      this.exocom.listen(exocomPort)
      done()
    }))
  })


  this.Given(/^an instance of this service$/, function(done) {
    portReservation.getPort(N( servicePort => {
      this.servicePort = servicePort
      this.exocom.registerService({name: serviceConfig.name, port: this.servicePort})
      this.process = new ExoService({ serviceName: serviceConfig.name,
                                       exocomPort: this.exocom.pullSocketPort,
                                       exorelayPort: servicePort })
      this.process.listen()
      this.process.on('online', function() { done() })
    }))
  })


  this.When(/^receiving the "([^"]*)" command$/, function(commandName) {
    this.exocom.reset()
    this.exocom.send({ service: serviceConfig.name,
                              name: commandName })
  })


  this.Then(/^this service replies with a "([^"]*)" message/, function(expectedMessageName, done) {
    this.exocom.onReceive( () => {
      const receivedMessages = this.exocom.receivedMessages
      expect(receivedMessages).to.have.length(1)
      expect(receivedMessages[0].name).to.equal(expectedMessageName)
      done()
    })
  })

}
