require! {
  'chalk' : {red}
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
    @service-config = if @config then yaml.safe-load fs.read-file-sync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    | !@_docker-file-exists!  =>  cp path.join(global.templates-path, 'docker', 'Dockerfile'), path.join(@config.root, 'Dockerfile')

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
