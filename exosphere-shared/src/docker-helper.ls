require! {
  \child_process
  \wait : {wait}
}


class DockerHelper

  @container-exists = (container) ->
    child_process.exec-sync('docker ps -a --format {{.Names}}') |> (.to-string!) |> (.split '\n') |> (.includes container)


  @container-is-running = (container-name) ->
    child_process.exec-sync('docker ps --format {{.Names}}/{{.Status}}') |> (.to-string!) |> (.split '\n') |> (.includes "#{container-name}/Up")


  @ensure-container-is-running = (container, done) ~>
    | @container-is-running container.container-name  =>  return done!
    | @container-exists container.container-name      =>  @start-container container.container-name; done!
    | otherwise                                       =>  @run-image container; wait 100, done


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


  @run-image = (container) ~>
    if container.container-name is \test-mongo
      child_process.exec-sync "docker run -d --name=#{container.container-name} -p 27017:27017 #{container.image}"
    else if container.container-name.includes \mongo
      child_process.exec-sync "docker run -d -v ~/Desktop/data:/data/db --name=#{container.container-name} #{container.image}"


  @start-container = (container-name) ~>
    child_process.exec-sync("docker start #{container-name}")


  @image-exists = (image) ->
    child_process.exec-sync("docker images #{image.author}/#{image.name}#{if image.version then ":#{image.version}" else ""}", "utf-8") |> (.includes "#{image.author}/#{image.name}")


  @remove-all-containers = ->
    all-containers = child_process.exec-sync 'docker ps -aq' |> (.to-string!)
    if all-containers
      child_process.exec-sync 'docker rm -f $(docker ps -aq)'



module.exports = DockerHelper
