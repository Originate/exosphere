require! {
  'exocom-mock' : ExoComMock
  'exoservice' : ExoService
  'chai' : {expect}
  'fs'
  'nitroglycerin' : N
  'port-reservation'
  'js-yaml' : yaml
}


service-config = yaml.safe-load(fs.read-file-sync('service.yml'), 'utf8')


module.exports = ->

  @Given /^an ExoCom server$/, (done) ->
    port-reservation.get-port N (exocom-port) ~>
      @exocom-port = exocom-port
      @exocom = new ExoComMock()
      @exocom.listen(exocom-port)
      done!


  @Given /^an instance of this service$/, (done) ->
    port-reservation.get-port N (service-port) ~>
      @service-port = service-port
      @exocom.register-service name: service-config.name, port: @service-port
      @process = new ExoService service-name: service-config.name, exocom-port: @exocom.pull-socket-port, exorelay-port: service-port
      @process.listen!
      @process.on 'online', -> done!


  @When /^receiving the "([^"]*)" command$/, (commandName) ->
    @exocom.reset!
    @exocom.send service: service-config.name, name: command-name


  @Then /^this service replies with a "([^"]*)" message/, (expectedMessageName, done) ->
    @exocom.on-receive  ~>
      received-messages = @.exocom.received-messages
      expect(received-messages).to.have.length 1
      expect(received-messages[0].name).to.equal expected-message-name
      done!
