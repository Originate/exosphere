require! {
  'chalk' : {red}
  'child_process'
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
  'prelude-ls' : {last}
  'require-yaml'
  'wait' : {wait-until}
}


# Runs a docker image
class DockerRunner extends EventEmitter

  (@name, @docker-config) ->


  start-service: ->
    @ensure-image-exists!
    @_remove-container!

    switch @name
      | \exocom    => @_run-container!
      | otherwise  =>
        wait-until (~> @get-docker-host 'exocom'), 10, ~>
          @docker-config.env.EXOCOM_HOST = @get-docker-host 'exocom'
          @_run-container!


  container-exists: (container) ->
   child_process.exec-sync("docker ps -a --format {{.Names}}", "utf8") |> (.includes container)


  ensure-image-exists: ->
    unless child_process.exec-sync("docker images #{@docker-config.author}/#{@docker-config.image}", "utf-8") |> (.includes "#{@docker-config.author}/#{@docker-config.image}")
      @emit 'error', "No Docker image exists for service #{@name}. Please run exo-setup."


  # Returns the IP address for the service with the given name
  get-docker-host: (container) ->
    child_process.exec-sync("docker inspect --format '{{ .NetworkSettings.IPAddress }}' #{container}", "utf8") if @container-exists container


  write: (text) ~>
    @emit 'output', {@name, text, trim: yes}


  _create-run-command: ->
    command = "
      docker run 
        --name=#{@docker-config.env.SERVICE_NAME} "
    for name, val of @docker-config.env
      command += " -e #{name}=#{val}"
    for name, port of @docker-config.publish
      command += " --publish #{port}"
    for link, container of  @docker-config.link
      command += " --link #{container}"
    command += " 
      #{@docker-config.author}/#{@docker-config.image} 
      #{@docker-config.start-command}"


  _on-container-error: ~>
    @emit 'error', "Service '#{@name}' crashed, shutting down application"


  # Removes a container with the given name from docker
  _remove-container: ->
    child_process.exec-sync "docker rm -f #{@name}" if @container-exists @name


  _run-container: ~>
    new ObservableProcess(@_create-run-command!,
                          stdout: {@write},
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code > 0  =>  @_on-container-error!
      ..wait @docker-config.start-text, ~>
        @emit 'online', @name



module.exports = DockerRunner
