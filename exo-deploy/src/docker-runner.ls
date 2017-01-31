require! {
  './docker-hub' : DockerHub
  '../../exosphere-shared' : {DockerHelper}
  'observable-process' : ObservableProcess
  'path'
}


class Docker

  (@app-config, @logger) ->
    process.env.AWS_ACCESS_KEY_ID ? throw new Error "AWS_ACCESS_KEY_ID not provided"
    process.env.AWS_SECRET_ACCESS_KEY ? throw new Error "AWS_SECRET_ACCESS_KEY not provided"
    {@version} = require '../../package.json'

  dockerhub-push: (done) ->
    new DockerHub @app-config, @logger
      ..push (err) -> done err


  start: (command-flag) ->
    image =
      author: 'originate'
      name: 'exo-deploy'
      version: @version

    if DockerHelper.image-exists image
      then @_run command-flag
      else
        @logger.log role: 'exo-deploy', text: "pulling ExoDeploy image version #{@version}"
        new ObservableProcess(DockerHelper.get-pull-command image,
                              stdout: {@write}
                              stderr: {@write})
          ..on 'ended', (exit-code) ~>
            | exit-code => return new Error "docker image originate/exo-deploy could not be pulled"
            @_run command-flag

  _run: (command-flag) ~>
    flags = "-v #{process.cwd!}:/var/app:ro " +
            "--env AWS_ACCESS_KEY_ID=#{process.env.AWS_ACCESS_KEY_ID} " +
            "--env AWS_SECRET_ACCESS_KEY=#{process.env.AWS_SECRET_ACCESS_KEY} "
    if process.env.MONGODB_USER then flags += "--env MONGODB_USER=#{process.env.MONGODB_USER} "
    if process.env.MONGODB_PW   then flags += "--env MONGODB_PW=#{process.env.MONGODB_PW}"
    new ObservableProcess("docker run #{flags} originate/exo-deploy:#{@version} #{command-flag}",
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code => return new Error "docker image originate/exo-deploy could not be run"


  write: (text) ~>
    @logger.log {role: 'exo-deploy', text, trim: no}


module.exports = Docker
