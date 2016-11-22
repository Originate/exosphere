require! {
  'chalk' : {red}
  'child_process'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  'exosphere-shared' : {call-args}
  'fs'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
  'prelude-ls' : {last}
  'require-yaml'
}


class ServiceRunner extends EventEmitter

  ({@name, @config, @logger}) ->
    @service-config = require path.join(@config.root, 'service.yml')


  start: (done) ~>
    @docker-config =
      author: @service-config.author
      image: path.basename @config.root
      start-command: @service-config.startup.command
      start-text: @service-config.startup['online-text']
      cwd: @config.root
      env:
        EXOCOM_PORT: @config.EXOCOM_PORT
        SERVICE_NAME: @name
      publish: @service-config.docker.publish if @service-config.docker
      link: @service-config.docker.link if @service-config.docker

    new DockerRunner {@name, @docker-config, @logger}
        ..start-service!
        ..on 'online', done
        ..on 'error', (message) ~> @emit 'error', error-message: message



module.exports = ServiceRunner
