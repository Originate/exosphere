const { defineSupportCode } = require('cucumber')
const ExoComMock = require('exocom-mock')
const { expect } = require('chai')
const fs = require('fs')
const N = require('nitroglycerin')
const portReservation = require('port-reservation')
const yaml = require('js-yaml')
const ObservableProcess = require('observable-process')


const serviceConfig = yaml.safeLoad(fs.readFileSync('service.yml'), 'utf8')


defineSupportCode(({ Given, When, Then }) => {
  Given(/^an ExoCom server$/, function(done) {
    portReservation.getPort(
      N(exocomPort => {
        this.exocomPort = exocomPort
        this.exocom = new ExoComMock()
        this.exocom.listen(exocomPort, done)
      })
    )
  })

  Given(/^an instance of this service$/, { timeout: 20 * 1000 }, function(
    done
  ) {
    const command = serviceConfig.development.scripts.run
    this.process = new ObservableProcess(command, {
      env: {
        EXOCOM_PORT: this.exocomPort,
        EXOCOM_HOST: 'localhost',
        ROLE: serviceConfig.type,
      },
      stdout: false,
      stderr: false,
    })
    this.process.wait(serviceConfig.startup['online-text'], done, 15 * 1000)
  })

  When(/^receiving the "([^"]*)" command$/, function(commandName) {
    this.exocom.reset()
    this.exocom.send({
      service: serviceConfig.type,
      name: commandName
    })
  })

  Then(/^this service replies with a "([^"]*)" message/, function(expectedMessageName, done) {
    this.exocom.onReceive( () => {
      const receivedMessages = this.exocom.receivedMessages
      expect(receivedMessages).to.have.length(1)
      expect(receivedMessages[0].name).to.equal(expectedMessageName)
      done()
    })
  })

})
