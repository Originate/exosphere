require! {
  'exosphere-shared' : {DockerHelper, Logger}
  'observable-process' : ObservableProcess
  'path'
}


class Docker

  ->
    @logger = new Logger
    process.env.AWS_ACCESS_KEY_ID ? throw new Error "AWS_ACCESS_KEY_ID not provided"
    process.env.AWS_SECRET_ACCESS_KEY ? throw new Error "AWS_SECRET_ACCESS_KEY not provided"
    @version = require path.join(__dirname, '../../package.json') |> (.version)


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
                              stdout: @_write
                              stderr: @_write)
          ..on 'ended', (exit-code) ~>
            | exit-code => return new Error "docker image originate/exo-deploy could not be pulled"
            @_run!

  _run: ~>
    flags = "-v #{process.cwd!}:/var/app:ro " +
            "--env AWS_ACCESS_KEY_ID=#{process.env.AWS_ACCESS_KEY_ID} " +
            "--env AWS_SECRET_ACCESS_KEY=#{process.env.AWS_SECRET_ACCESS_KEY} "
    new ObservableProcess("docker run #{flags} originate/exo-deploy:#{@version}",
                          stdout: @_write
                          stderr: @_write)
      ..on 'ended', (exit-code) ~>
        | exit-code => return new Error "docker image originate/exo-deploy could not be run"


  _write: (text) ->
    @logger.log {'exo-deploy', text, trim: no}


module.exports = Docker
