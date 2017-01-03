require! {
  'chalk' : {red}
  'child_process'
  'dashify'
  'events' : {EventEmitter}
  '../../exosphere-shared' : {templates-path, call-args, DockerHelper}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
  'shelljs' : {cp}
}


class DockerSetup extends EventEmitter

  ({@name, @logger, @config}) ->
    # Services being pulled from Docker Hub will not have a '@service-config'
    try
      @service-config = if @config then yaml.safe-load fs.read-file-sync(path.join(@config.root, 'service.yml'), 'utf8')
    catch
      return

  start: (done) ~>
    | !@service-config  =>
      if @config.docker-image
        image = @config.docker-image |> (.split '/')
        new ObservableProcess((DockerHelper.get-pull-command author: image[0], name: image[1]),
                              stdout: {@write}
                              stderr: {@write})
          ..on 'ended', (exit-code) ~> done!
        return
      else
        @logger.log "No location or docker-image specified"
        process.exit 1
    | !@_docker-file-exists!    =>  cp path.join(templates-path, 'docker', 'Dockerfile'), path.join(@config.root, 'Dockerfile')

    @logger.log name: @name, text: "preparing Docker image"
    @_build-docker-image done


  _build-docker-image: (done) ~>
    new ObservableProcess(call-args(DockerHelper.get-build-command author: @service-config.author, name: dashify(@service-config.title)),
                          cwd: @config.root
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>  @logger.log name: @name, text: "Docker setup finished"
        | otherwise       =>
          @logger.log name: @name, text: "Docker setup failed"
          process.exit exit-code
        done!


  _docker-file-exists: ~>
    try
      fs.exists-sync path.join(@config.root, 'Dockerfile')
    catch
      no


  write: (text) ~>
    @emit 'output', {@name, text, trim: yes}



module.exports = DockerSetup
