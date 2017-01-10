require! {
  'chalk' : {red}
  'child_process'
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper}
  'fs'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
  'prelude-ls' : {last}
  'wait' : {wait-until}
}


# Runs a docker image
class DockerRunner extends EventEmitter

  ({@role, @docker-config, @logger}) ->


  start-service: ->
    unless DockerHelper.image-exists author: @docker-config.author, name: @docker-config.image
      return @emit 'error', "No Docker image exists for service '#{@role}'. Please run exo-setup."
    DockerHelper.remove-container @role

    switch @role
      | \exocom    => @_run-container!
      | otherwise  =>
        wait-until (~> DockerHelper.get-docker-ip 'exocom'), 10, ~>
          @docker-config.env.EXOCOM_HOST = DockerHelper.get-docker-ip 'exocom'
          @_check-dependency-containers!
          @_run-container!


  write: (text) ~>
    @logger.log {name: @role, text, trim: yes}


  _create-run-command: ->
    command = "
      docker run 
        --name=#{@docker-config.env.ROLE} "
    for name, val of @docker-config.env
      command += " -e #{name}=#{val}"
    for name, port of @docker-config.publish
      command += " --publish #{port}"
    command += " 
      #{@docker-config.author}/#{@docker-config.image} 
      #{@docker-config.start-command}"


  _on-container-error: ~>
    @emit 'error', "Service '#{@role}' crashed, shutting down application"


  _run-container: ~>
    @docker-container = new ObservableProcess(@_create-run-command!,
                                              stdout: {@write},
                                              stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code > 0 and not killed   =>  @_on-container-error!
      ..wait @docker-config.start-text, ~>
        @logger.log name: 'exo-run', text: "'#{@role}' is running"
        @emit 'online'


  _check-dependency-containers: ~>
    for dependency of @docker-config.dependencies
      app-dependency = "#{@docker-config.app-name}-#{dependency}"
      DockerHelper.ensure-container-is-running app-dependency, dependency
      @docker-config.env[dependency.to-upper-case!] = DockerHelper.get-docker-ip app-dependency



module.exports = DockerRunner
