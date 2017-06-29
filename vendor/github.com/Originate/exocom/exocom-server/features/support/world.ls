require! {
  'async'
  'chai' : {expect}
  'dim-console'
  'jsdiff-console'
  './mock-service': MockService
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'record-http' : HttpRecorder
  'request'
  './text-tools' : {ascii}
  'wait' : {wait-until, wait}
}


observableProcessOptions = if process.env.DEBUG_EXOCOM_SERVER
  stdout: dim-console.process.stdout
  stderr: dim-console.process.stderr
else
  stdout: no
  stderr: no


# Provides steps for end-to-end testing of the service as a stand-alone binary
World = !->

  @create-exocom-instance = ({port, service-routes = '{}'}, done) ->
    env =
      PORT : port
      SERVICE_ROUTES : service-routes
    @port = port
    @process = new ObservableProcess "bin/exocom", stdout: observableProcessOptions.stdout, stderr: observableProcessOptions.stderr, env: env
      ..wait "WebSocket listener online at port #{@port}", done


  @create-mock-service-at-port = ({client-name, port, namespace}, done) ->
    (@service-mocks or= {})[client-name] = new MockService {port, client-name, namespace}
    @service-mocks[client-name].connect {}, ->
      wait 200, done


  @run-exocom-at-port = (port, _expect-error, done) ->
    env =
      PORT : port
    @process = new ObservableProcess "bin/exocom", stdout: observableProcessOptions.stdout, stderr: observableProcessOptions.stderr, env: env
    done!


  @service-sends-message = ({service-name, message-name}, done) ->
    request-data =
      sender: service-name
      payload: ''
      id: '123'
      name: message-name
    @service-mocks[service-name].send request-data
    done!


  @service-sends-reply = ({service-name, message-name, response-to}, done) ->
    request-data =
      sender: service-name
      payload: ''
      id: '123'
      response-to: response-to
      name: message-name
    @service-mocks[service-name].send request-data
    done!


  @verify-abort-with-notification = (text, done) ->
    @process.wait text, ~>
      wait-until (~> @process.ended), done


  @verify-exocom-broadcasted-message = ({message-name, sender, receivers}, done) ->
    @process.wait "#{sender} is broadcasting '#{message-name}' to the #{receivers.join ', '}", done


  @verify-exocom-signaled-string = (text, done) ->
    [...main-parts, response-time-msg] = text.split '  '
    @process.wait "#{main-parts.join '  '}  ", done


  @verify-exocom-received-message = (message-name, done) ->
    @process.wait "broadcasting '#{message-name}'", done


  @verify-exocom-received-reply = (message-name, done) ->
    @process.wait "broadcasting '#{message-name}'", done


  @verify-routing-setup = (expected-routing, done) ->
    request "http://localhost:#{@port}/config.json", (err, result, body) ->
      expect(err).to.be.null
      expect(result.status-code).to.equal 200
      jsdiff-console JSON.parse(body).routes, expected-routing, done


  @verify-listening-at-ports = (port, done) ->
    messages = []
    messages.push "WebSocket listener online at port #{port}" if port
    messages.push "HTTP service online at port #{port}" if port
    async.each messages,
               ((text, cb) ~> @process.wait text, cb),
               done


  @verify-sent-calls = ({client-name, message-name, id = '123', response-to}, done) ->
    service = @service-mocks[client-name]
    wait-until (~> service.received-messages[0]?.name is message-name), 1, ~>
      expected =
        name: message-name
        id: id
        payload: ''
      received-message = service.received-messages[0]
      expected.timestamp = that if received-message.timestamp
      expected.response-time = that if received-message.response-time
      expected.response-to = that if response-to
      jsdiff-console received-message, expected, done


  @verify-service-setup = (service-data, done) ->
    request "http://localhost:#{@port}/config.json", (err, result, body) ->
      expect(err).to.be.null
      expect(result.status-code).to.equal 200
      jsdiff-console JSON.parse(body).clients, service-data, done



module.exports = ->
  @World = World
