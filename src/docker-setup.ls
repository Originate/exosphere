require! {
  'chalk' : {red}
  'events' : {EventEmitter}
  'exosphere-shared' : {templates-path, call-args}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
  'prelude-ls' : {last}
  'require-yaml'
  'shelljs' : {cp}
}


class DockerSetup extends EventEmitter

  ({@name, @logger, @config}) ->
    @service-config = require path.join(@config.root, 'service.yml')


  start: (done) ~>
    | @_docker-file-exists!  =>  @logger.log name: @name, text: "Docker image already exists"; return done!
    @logger.log name: @name, text: "starting setup of Docker image"

    service-name = @config.root.split path.sep |> last
    cp path.join(templates-path, 'docker', 'Dockerfile'), path.join(@config.root, 'Dockerfile')
    new ObservableProcess(call-args("docker build --build-arg SERVICE_NAME=#{service-name} -t #{@service-config.author}/#{service-name} ."),
                          cwd: @config.root,
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
      fs.access-sync path.join(@config.root, 'Dockerfile')
    catch
      no

  write: (text) ~>
    @emit 'output', {@name, text}



module.exports = DockerSetup
