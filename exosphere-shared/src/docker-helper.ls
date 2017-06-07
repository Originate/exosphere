require! {
  'async'
  \prelude-ls : {any, head, map}
  'dockerode' : Docker
  'text-stream-search' : TextStreamSearch
}

docker = new Docker

# Helper class used to manage Docker processes not started by docker-compose
class DockerHelper

  @start-container = ({Image, name, HostConfig, online-text}, done) ->
    console.log Image
    DockerHelper.list-running-containers (err, running-containers) ->
      | err                              => done err
      | running-containers.includes name => done!
      | otherwise                        =>
        docker.create-container {Image, name, HostConfig}, (err, container) -> 
          | err => done err
          container.attach {stream: true, stdout: true, stederr: true}, (err, stream) ->
            text-stream-search = new TextStreamSearch stream    
            container.start (err) ->
              | err => done err
              text-stream-search.wait online-text, done


  @remove-container = ({name}, done) ->
    docker.list-containers {name}, (err, containers) ->
      | err => done err
      DockerHelper._force-remove-containers containers, done


  @remove-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err                => done err
      | !containers.length => done!
      | otherwise          => DockerHelper._force-remove-containers containers, done


  @list-running-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err => done err
      # Names field is printed like: Names: [ '/exocom' ]
      done null, map((.Names?[0] |> (.replace '/', '')), containers)


  @list-images = (done) ->
    docker.list-images (err, images) ->
      | err => done err
      # Image name is printed like: RepoTags: [ 'exocom:latest' ]
      done null, map((.RepoTags |> head |> (.split ':') |> head), images)


  @get-dangling-images = (done) ->
    docker.list-images {"filters": '{"dangling": ["true"]}'}, (err, images) ->
      | err => done err
      done null, images


  @force-remove-images = (images, done) ->
    if images.length == 0 then done!
    else
      async.map-series images, ((image, done) -> docker.get-image(image.Id).remove {force:true}, ((err) -> done err)), (err) -> 
        | err => done err
        done!


  @get-dangling-volumes = (done) ->
    docker.list-volumes {"filters": '{"dangling": ["true"]}'}, (err, volumes) ->
      | err => done err
      done null, (volumes.Volumes or [])


  @force-remove-volumes = (volumes, done) ->
    if volumes.length == 0 then done!
    else
      async.map-series volumes, ((volume, done) -> docker.get-volume(volume.Name).remove {force:true}, ((err) -> done err)), (err) -> 
        | err => done err
        done!


  @_force-remove-containers = (containers, done) ->
    for container in containers
      docker.get-container(container.Id).remove {force:true}, (err) ->
        | err => done err
        done!


module.exports = DockerHelper
