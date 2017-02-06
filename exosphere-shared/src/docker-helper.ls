require! {
  \child_process
  \observable-process : ObservableProcess
}


class DockerHelper

  @container-exists = (container) ->
    child_process.exec-sync('docker ps -a --format {{.Names}}') |> (.to-string!) |> (.split '\n') |> (.includes container)


  @container-is-running = (container-name) ->
    child_process.exec-sync('docker ps --format {{.Names}}/{{.Status}}') |> (.to-string!) |> (.split '\n') |> (.includes "#{container-name}/Up")


  @ensure-container-is-running = (container, done) ~>
    | @container-is-running container.container-name  =>  return done!
    | @container-exists container.container-name      =>  @start-container container, (err) ~> done err
    | otherwise                                       =>  @run-image container, (err) ~> done err


  @get-build-command = (image, build-flags) ->
    return "docker build -t #{image.author}/#{image.name} #{if build-flags then build-flags else ""} ."


  @get-config = (image) ->
    child_process.exec-sync("docker run --rm=true #{image} cat service.yml", 'utf8') |> (.to-string!)


  @get-docker-ip = (container) ->
    child_process.exec-sync("docker inspect --format '{{ .NetworkSettings.IPAddress }}' #{container}", "utf8") if @container-exists container


  @get-docker-images = ->
    child_process.exec-sync 'docker images'


  @get-pull-command = (image) ->
    return "docker pull #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}"


  @remove-container = (container) ->
    child_process.exec-sync "docker rm -f #{container}" if @container-exists container


  @run-image = (container, done) ~>
      #TODO: Clean up the flags to the docker run command
      new ObservableProcess("docker run #{if container.volume then that else ''} #{if container.port then that else ''} --name=#{container.container-name} #{container.dependency-name}",
                            stdout: false,
                            stderr: false)
        ..on 'ended', (exit-code, killed) ->
          | exit-code > 0 and not killed  =>  return done "Dependency #{container.container-name} failed to run, shutting down"
        ..wait container.online-text, done


  @start-container = (container, done) ~>
    new ObservableProcess("docker start -a #{container.container-name}",
                            stdout: false,
                            stderr: false)
      ..on 'ended', (exit-code, killed) ->
        | exit-code > 0 and not killed  =>  return done "Dependency #{container.container-name} failed to start, shutting down"
      ..wait container.online-text, done


  @image-exists = (image) ->
    child_process.exec-sync("docker images #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}", "utf-8") |> (.includes "#{image.author}/#{image.name}")


  @remove-all-containers = ->
    all-containers = child_process.exec-sync 'docker ps -aq' |> (.to-string!)
    if all-containers
      child_process.exec-sync 'docker rm -f $(docker ps -aq)'



module.exports = DockerHelper
