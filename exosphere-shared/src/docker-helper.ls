require! {
  \child_process
}


class DockerHelper

  @container-exists = (container) ->
    child_process.exec-sync('docker ps -a --format {{.Names}}') |> (.to-string!) |> (.split '\n') |> (.includes container)


  @container-is-running = (container-name) ->
    child_process.exec-sync('docker ps --format {{.Names}}/{{.Status}}') |> (.to-string!) |> (.split '\n') |> (.includes "#{container-name}/Up")


  @ensure-container-is-running = (container-name, image) ->
    | @container-is-running container-name  =>  return
    | @container-exists container-name      =>  @start-container container-name
    | otherwise                             =>  @run-image container-name, image


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


  @run-image = (container, image) ->
    if container is \test-mongo
      child_process.exec-sync "docker run -d --name=#{container} -p 27017:27017 #{image}"
    else
      child_process.exec-sync "docker run -d --name=#{container} #{image}"


  @start-container = (container-name) ->
    child_process.exec-sync("docker start #{container-name}") if @container-exists container-name


  @image-exists = (image) ->
    child_process.exec-sync("docker images #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}", "utf-8") |> (.includes "#{image.author}/#{image.name}")


  @remove-all-containers = ->
    all-containers = child_process.exec-sync 'docker ps -aq' |> (.to-string!)
    if all-containers
      child_process.exec-sync 'docker rm -f $(docker ps -aq)'



module.exports = DockerHelper
