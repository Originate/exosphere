require! {
  \child_process
}


class DockerHelper

  @container-exists = (container) ->
    child_process.exec-sync('docker ps -a --format {{.Names}}', 'utf8') |> (.includes container)


  @ensure-container-is-running = (container, image) ->
    return if child_process.exec-sync('docker ps --format {{.Names}}/{{.Status}}' , 'utf8') |> (.includes "#{container}/Up")
    @remove-container that if @container-exists container
    @run-image container, image


  @get-build-command = (image, build-flags) ->
    return "docker build -t #{image.author}/#{image.name} #{if build-flags then build-flags else ""} ."


  @get-config = (image) ->
    return child_process.exec-sync("docker run #{image} cat service.yml", 'utf8') |> (.to-string!)


  @get-docker-ip = (container) ->
    child_process.exec-sync("docker inspect --format '{{ .NetworkSettings.IPAddress }}' #{container}", "utf8") if @container-exists container


  @get-docker-images = ->
    child_process.exec-sync 'docker images'


  @get-pull-command = (image) ->
    return "docker pull #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}"


  @remove-container = (container) ->
    child_process.exec-sync "docker rm -f #{container}" if @container-exists container


  @run-image = (container, image) ->
    child_process.exec-sync "docker run -d --name=#{container} #{image}"

  @image-exists = (image) ->
    child_process.exec-sync("docker images #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}", "utf-8") |> (.includes "#{image.author}/#{image.name}")



module.exports = DockerHelper
