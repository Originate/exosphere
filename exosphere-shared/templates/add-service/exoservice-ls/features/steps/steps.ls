require! {
  'cucumber': {defineSupportCode}
  'exocom-mock' : ExoComMock
  'exoservice' : ExoService
  'chai' : {expect}
  'fs'
  'nitroglycerin' : N
  'port-reservation'
  'js-yaml' : yaml
  'wait': {wait}
}


service-config = yaml.safe-load(fs.read-file-sync('service.yml'), 'utf8')


defineSupportCode ({Given, When, Then}) ->

  Given /^an ExoCom server$/, (done) ->
    port-reservation.get-port N (exocom-port) ~>
      @exocom-port = exocom-port
      @exocom = new ExoComMock()
      @exocom.listen exocom-port, done


  Given /^an instance of this service$/, (done) ->
    @process = new ExoService role: service-config.type, exocom-port: @exocom-port, exocom-host: 'localhost'
      ..connect!
      ..on 'online', ~> wait 10, done


  When /^receiving the "([^"]*)" command$/, (commandName) ->
    @exocom.reset!
    @exocom.send service: service-config.type, name: command-name


  Then /^this service replies with a "([^"]*)" message/, (expectedMessageName, done) ->
    @exocom.on-receive  ~>
      received-messages = @.exocom.received-messages
      expect(received-messages).to.have.length 1
      expect(received-messages[0].name).to.equal expected-message-name
      done!
