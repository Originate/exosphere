require! {
  './docker-hub' : DockerHub
  'exosphere-shared' : {DockerHelper}
  'observable-process' : ObservableProcess
  'path'
}


class Docker

  (@app-config, @logger) ->
    process.env.AWS_ACCESS_KEY_ID ? throw new Error "AWS_ACCESS_KEY_ID not provided"
    process.env.AWS_SECRET_ACCESS_KEY ? throw new Error "AWS_SECRET_ACCESS_KEY not provided"
    @version = require path.join(__dirname, '../../package.json') |> (.version)


  dockerhub-push: (done) ->
    new DockerHub @app-config, @logger
      ..push (err) -> done err


  start: ->
    image =
      author: 'originate'
      name: 'exo-deploy'
      version: @version

    if DockerHelper.image-exists image
      then @_run!
      else
        @logger.log name: 'exo-deploy', text: "pulling ExoDeploy image version #{@version}"
        new ObservableProcess(DockerHelper.get-pull-command image,
                              stdout: {@write}
                              stderr: {@write})
          ..on 'ended', (exit-code) ~>
            | exit-code => return new Error "docker image originate/exo-deploy could not be pulled"
            @_run!

  _run: ~>
    flags = "-v #{process.cwd!}:/var/app:ro " +
            "--env AWS_ACCESS_KEY_ID=#{process.env.AWS_ACCESS_KEY_ID} " +
            "--env AWS_SECRET_ACCESS_KEY=#{process.env.AWS_SECRET_ACCESS_KEY} "
    new ObservableProcess("docker run #{flags} originate/exo-deploy:#{@version}",
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code => return new Error "docker image originate/exo-deploy could not be run"


  write: (text) ~>
    @logger.log {name: 'exo-deploy', text, trim: no}


module.exports = Docker
