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

  ({@role, @logger}) ->
    @service-config = if @config then yaml.safe-load fs.read-file-sync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    | !@service-config        =>  return @_setup-external-service done
    | !@_docker-file-exists!  =>  cp path.join(templates-path, 'docker', 'Dockerfile'), path.join(@config.root, 'Dockerfile')

    @logger.log {@role, text: "preparing Docker image"}
    @_build-docker-image done


  _build-docker-image: (done) ~>
    new ObservableProcess(call-args(DockerHelper.get-build-command author: @service-config.author, name: dashify(@service-config.type))
                          cwd: @config.root
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>  @logger.log {@role, text: "Docker setup finished"}
        | otherwise       =>
          @logger.log {@role, text: "Docker setup failed"}
          process.exit exit-code
        done!


  _docker-file-exists: ~>
    try
      fs.exists-sync path.join(@config.root, 'Dockerfile')
    catch
      no


  _setup-external-service: (done) ~>
    throw new Error red "No location or docker-image specified" unless @config.docker-image
    image = @config.docker-image |> (.split '/')
    new ObservableProcess((DockerHelper.get-pull-command author: image[0], name: image[1]),
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code isnt 0  =>
          @logger.log {@role, text: 'Docker setup failed'}
          process.exit exit-code
        | _                 =>  done!


  write: (text) ~>
    @logger.log {@role, text, trim: yes}



module.exports = DockerSetup
