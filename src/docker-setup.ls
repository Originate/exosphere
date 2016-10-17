require! {
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

  (@name, @config) ->
    @service-config = require path.join(@config.root, 'service.yml')


  start: (done) ~>
    | @_docker-file-exists!  =>  @emit 'docker-exists', @name; return done!
    @emit 'docker-start', @name

    service-name = @config.root.split path.sep |> last
    cp path.join(templates-path, 'docker', 'Dockerfile'), path.join(@config.root, 'Dockerfile')
    new ObservableProcess(call-args("docker build --build-arg SERVICE_NAME=#{service-name} -t #{@service-config.author}/#{service-name} ."),
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>  @emit 'docker-finished', @name
        | otherwise       =>  @emit 'error', @name, exit-code
        done!


  _docker-file-exists: ~>
    try
      fs.access-sync path.join(@config.root, 'Dockerfile')
    catch
      no

  write: (text) ~>
    @emit 'output', {@name, text}



module.exports = DockerSetup
